"""YANG test plugin for correctness of "pattern" statement regexes.
"""

from io import StringIO
import optparse
import re

from pyang import plugin
from pyang import statements
from pyang import error
from pyang import grammar
from pyang.error import err_add

from lxml import etree

def pyang_plugin_init():
    plugin.register_plugin(PatternTestPlugin())

class PatternTestPlugin(plugin.PyangPlugin):
    def __init__(self):
        plugin.PyangPlugin.__init__(self)

    def post_validate_ctx(self, ctx, modules):
        if not ctx.opts.check_patterns:
            return
        check_patterns(ctx, modules)

    def add_opts(self, optparser):
        optlist = [
            optparse.make_option("--check-patterns",
                                 dest="check_patterns",
                                 action="store_true",
                                 help="Validate pattern-test-pass and"
                                 "pattern-test-fail cases for string"
                                 "restrictions."),
        ]
        optparser.add_options(optlist)

    def setup_ctx(self, ctx):
        if not ctx.opts.check_patterns:
            return
        # Test case failure states.
        error.add_error_code(
            'VALID_PATTERN_DOES_NOT_MATCH', 2,
            'valid pattern "%s" does not match typedef "%s"')
        error.add_error_code(
            'INVALID_PATTERN_MATCH', 2,
            'invalid pattern "%s" matches typedef "%s"')

        # Error states.
        error.add_error_code(
            'NO_TEST_PATTERNS', 1,
            'leaf "%s" does not have any test cases')
        error.add_error_code(
            'UNRESTRICTED_TYPE', 1,
            'string type "%s" is unrestricted')

typedef_usage_stmt_regex = re.compile(r'([^\s:]+:)?([^\s:]+)')

debug = False

def is_statement_pass_testcase(statement):
    try:
        return statement.keyword[1] == 'pattern-test-pass'
    except AttributeError:
        return False

def is_statement_fail_testcase(statement):
    try:
        return statement.keyword[1] == 'pattern-test-fail'
    except AttributeError:
        return False

def dnf_patterns(ctx, typestmt, prefix_to_mod_name):
    """dnf_patterns obtains a DNF of the type's patterns.

    The return value is a list of list of patterns representing the OR of ANDs,
    aka. sum-of-products (DNF), of the type's patterns.

    Satisfying any product group (minterm) of the DNF means that the pattern is
    satisfied.

    e.g.
    typedef ip-version {
        union {
            type ip-version4;
            type ip-version6;
        }
    }

    typedef ip-version4 {
        type string {
            pattern 'ip[vV]4';
            pattern '\w+';
        }
    }

    typedef ip-version6 {
        type string {
            pattern 'ip[vV]6';
            pattern '\w+';
        }
    }

    output: [["ip[vV]4", "\w+"],  ["ip[vV6]", "\w+"]]
    """
    patterns = []
    dnf_patterns_aux(ctx, typestmt, prefix_to_mod_name, patterns)
    return patterns

def dnf_patterns_aux(ctx, typestmt, prefix_to_mod_name, patterns):
    """dnf_patterns_aux is the recursive pattern finder for dnf_patterns.

    Recursively find the base string types of the given type, and append its
    pattern statements into the given patterns list. The "pattern" statement
    specification + YANG type hierarchy naturally makes each string type a
    minterm in the DNF, no matter how deeply nested it is inside a union.
    """
    if typestmt.keyword != "type":
        return

    # TODO(wenovus): Support patterns found in derived string types.
    # See https://tools.ietf.org/html/rfc7950#section-9.4.5
    if typestmt.arg == "string":
        new_patterns = typestmt.search("pattern")
        if not new_patterns:
            # If we have a string type without restrictions, then there is no
            # reason to have a test for the type.
            err_add(ctx.errors, typestmt.pos, 'UNRESTRICTED_TYPE', (typestmt.arg))
        else:
            # Each string type will form a minterm for the DNF.
            patterns.append([pat.arg for pat in new_patterns])
        return

    if typestmt.arg == "union":
        # Each string union subtype will form a minterm for the DNF.
        for substmt in typestmt.search("type"):
            dnf_patterns_aux(ctx, substmt, prefix_to_mod_name, patterns)
        return

    # If we encounter a derived type, descend into it to find its constituent
    # minterms.
    base_typestmt = typestmt_from_derived_typestmt(ctx, typestmt, prefix_to_mod_name)
    if base_typestmt:
        dnf_patterns_aux(ctx, base_typestmt, prefix_to_mod_name, patterns)

def typestmt_from_derived_typestmt(ctx, typestmt, prefix_to_mod_name):
    """typestmt_from_derived_typestmt retrieves the type statement from the
    input typedef usage statement.

    e.g.
    Module A
    leaf foo {
        type ip-types:ipv4;   <-- input statement
    }

    Module B
    typedef ipv4 {
        type string {         <-- output statement
            pattern 'ip[vV]4';
        }
    }
    """
    match = typedef_usage_stmt_regex.match(typestmt.arg)
    if not match:
        return
    typemod = None
    if not match.group(1):
        typemod = ctx.get_module(typestmt.main_module().arg)
    elif match.group(1).endswith(":") and match.group(1)[:-1] in prefix_to_mod_name:
        typemod = ctx.get_module(prefix_to_mod_name[match.group(1)[:-1]])
    if not typemod:
        return

    for typedef in typemod.search("typedef"):
        if typedef.arg == match.group(2):
            return typedef.search_one("type")

def check_patterns(ctx, mods):
    """check_patterns executes any pattern tests in the YANG file."""

    xsd_doc_str = """<xsd:schema xmlns:xsd="http://www.w3.org/2001/XMLSchema">
    <xsd:element name="val">
        <xsd:simpleType>
            <xsd:restriction base="xsd:string">
                <xsd:pattern value="{}"/>
            </xsd:restriction>
        </xsd:simpleType>
    </xsd:element>
</xsd:schema>"""

    for mod in mods:
        # Build map of {import prefix: module name}
        prefix_to_mod_name = {}
        for s in mod.search("import"):
            prefix_to_mod_name[s.search_one("prefix").arg] = s.arg

        # LIMITATION: Only direct leaves can currently be tested.
        for leaf in mod.search("leaf"):
            typestmt = leaf.search_one("type")
            pattern_lists = dnf_patterns(ctx, typestmt, prefix_to_mod_name)
            if not pattern_lists:
                err_add(ctx.errors, typestmt.pos, 'UNRESTRICTED_TYPE',
                        (typestmt.arg))
                continue

            has_at_least_one_test = False
            for s in leaf.substmts:
                xml_doc = etree.parse(StringIO("<val>{}</val>".format(s.arg)))
                if is_statement_pass_testcase(s):
                    # Check whether any minterm in the DNF is satisfied.
                    has_at_least_one_test = True
                    for patterns in pattern_lists:
                        minterm_matches = True
                        for pattern in patterns:
                            f = StringIO(xsd_doc_str.format(pattern))
                            xsd_doc = etree.parse(f)
                            xsd = etree.XMLSchema(xsd_doc)
                            if not xsd.validate(xml_doc):
                                if debug:
                                    print("{} doesn't match {}".format(s.arg, pattern))
                                minterm_matches = False
                            elif debug:
                                print("{} matches".format(s.arg))
                        if minterm_matches:
                            break
                    else: # Exhausted all minterms.
                        err_add(ctx.errors, s.pos,
                                'VALID_PATTERN_DOES_NOT_MATCH', (s.arg,
                                                                 typestmt.arg))
                elif is_statement_fail_testcase(s):
                    # Check that no minterm in the DNF is satisfied.
                    has_at_least_one_test = True
                    for patterns in pattern_lists:
                        minterm_matches = True
                        for pattern in patterns:
                            f = StringIO(xsd_doc_str.format(pattern))
                            xsd_doc = etree.parse(f)
                            xsd = etree.XMLSchema(xsd_doc)
                            if not xsd.validate(xml_doc):
                                minterm_matches = False
                        if minterm_matches:
                            err_add(ctx.errors, s.pos, 'INVALID_PATTERN_MATCH',
                                    (s.arg, typestmt.arg))
                            break

            if not has_at_least_one_test:
                err_add(ctx.errors, leaf.pos, 'NO_TEST_PATTERNS', (leaf.arg))
