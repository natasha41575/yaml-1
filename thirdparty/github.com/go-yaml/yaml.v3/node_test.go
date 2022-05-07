//
// Copyright (c) 2011-2019 Canonical Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package yaml_test

import (
	"bytes"
	"fmt"
	"os"

	"io"
	"strings"

	. "gopkg.in/check.v1"
	. "sigs.k8s.io/yaml/thirdparty/github.com/go-yaml/yaml.v3"
)

var nodeTests = []struct {
	yaml string
	node Node
}{
	{
		"null\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   ScalarNode,
				Value:  "null",
				Tag:    "!!null",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"[encode]null\n",
		Node{},
	}, {
		"foo\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   ScalarNode,
				Value:  "foo",
				Tag:    "!!str",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"\"foo\"\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   ScalarNode,
				Style:  DoubleQuotedStyle,
				Value:  "foo",
				Tag:    "!!str",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"'foo'\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   ScalarNode,
				Style:  SingleQuotedStyle,
				Value:  "foo",
				Tag:    "!!str",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"!!str 123\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   ScalarNode,
				Style:  TaggedStyle,
				Value:  "123",
				Tag:    "!!str",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		// Although the node isn't TaggedStyle, dropping the tag would change the value.
		"[encode]!!binary gIGC\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   ScalarNode,
				Value:  "gIGC",
				Tag:    "!!binary",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		// Item doesn't have a tag, but needs to be binary encoded due to its content.
		"[encode]!!binary gIGC\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   ScalarNode,
				Value:  "\x80\x81\x82",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		// Same, but with strings we can just quote them.
		"[encode]\"123\"\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   ScalarNode,
				Value:  "123",
				Tag:    "!!str",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"!tag:something 123\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   ScalarNode,
				Style:  TaggedStyle,
				Value:  "123",
				Tag:    "!tag:something",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"[encode]!tag:something 123\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   ScalarNode,
				Value:  "123",
				Tag:    "!tag:something",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"!tag:something {}\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   MappingNode,
				Style:  TaggedStyle | FlowStyle,
				Tag:    "!tag:something",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"[encode]!tag:something {}\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   MappingNode,
				Style:  FlowStyle,
				Tag:    "!tag:something",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"!tag:something []\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   SequenceNode,
				Style:  TaggedStyle | FlowStyle,
				Tag:    "!tag:something",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"[encode]!tag:something []\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   SequenceNode,
				Style:  FlowStyle,
				Tag:    "!tag:something",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"''\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   ScalarNode,
				Style:  SingleQuotedStyle,
				Value:  "",
				Tag:    "!!str",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"|\n  foo\n  bar\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   ScalarNode,
				Style:  LiteralStyle,
				Value:  "foo\nbar\n",
				Tag:    "!!str",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"true\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   ScalarNode,
				Value:  "true",
				Tag:    "!!bool",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"-10\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   ScalarNode,
				Value:  "-10",
				Tag:    "!!int",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"4294967296\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   ScalarNode,
				Value:  "4294967296",
				Tag:    "!!int",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"0.1000\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   ScalarNode,
				Value:  "0.1000",
				Tag:    "!!float",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"-.inf\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   ScalarNode,
				Value:  "-.inf",
				Tag:    "!!float",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		".nan\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   ScalarNode,
				Value:  ".nan",
				Tag:    "!!float",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"{}\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   MappingNode,
				Style:  FlowStyle,
				Value:  "",
				Tag:    "!!map",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"a: b c\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   MappingNode,
				Value:  "",
				Tag:    "!!map",
				Line:   1,
				Column: 1,
				Content: []*Node{{
					Kind:   ScalarNode,
					Value:  "a",
					Tag:    "!!str",
					Line:   1,
					Column: 1,
				}, {
					Kind:   ScalarNode,
					Value:  "b c",
					Tag:    "!!str",
					Line:   1,
					Column: 4,
				}},
			}},
		},
	}, {
		"a:\n  b: c\n  d: e\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   MappingNode,
				Tag:    "!!map",
				Line:   1,
				Column: 1,
				Content: []*Node{{
					Kind:   ScalarNode,
					Value:  "a",
					Tag:    "!!str",
					Line:   1,
					Column: 1,
				}, {
					Kind:   MappingNode,
					Tag:    "!!map",
					Line:   2,
					Column: 3,
					Content: []*Node{{
						Kind:   ScalarNode,
						Value:  "b",
						Tag:    "!!str",
						Line:   2,
						Column: 3,
					}, {
						Kind:   ScalarNode,
						Value:  "c",
						Tag:    "!!str",
						Line:   2,
						Column: 6,
					}, {
						Kind:   ScalarNode,
						Value:  "d",
						Tag:    "!!str",
						Line:   3,
						Column: 3,
					}, {
						Kind:   ScalarNode,
						Value:  "e",
						Tag:    "!!str",
						Line:   3,
						Column: 6,
					}},
				}},
			}},
		},
	}, {
		"a:\n  - b: c\n    d: e\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   MappingNode,
				Tag:    "!!map",
				Line:   1,
				Column: 1,
				Content: []*Node{{
					Kind:   ScalarNode,
					Value:  "a",
					Tag:    "!!str",
					Line:   1,
					Column: 1,
				}, {
					Kind:   SequenceNode,
					Tag:    "!!seq",
					Line:   2,
					Column: 3,
					Content: []*Node{{
						Kind:   MappingNode,
						Tag:    "!!map",
						Line:   2,
						Column: 5,
						Content: []*Node{{
							Kind:   ScalarNode,
							Value:  "b",
							Tag:    "!!str",
							Line:   2,
							Column: 5,
						}, {
							Kind:   ScalarNode,
							Value:  "c",
							Tag:    "!!str",
							Line:   2,
							Column: 8,
						}, {
							Kind:   ScalarNode,
							Value:  "d",
							Tag:    "!!str",
							Line:   3,
							Column: 5,
						}, {
							Kind:   ScalarNode,
							Value:  "e",
							Tag:    "!!str",
							Line:   3,
							Column: 8,
						}},
					}},
				}},
			}},
		},
	}, {
		"a: # AI\n  - b\nc:\n  - d\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   MappingNode,
				Tag:    "!!map",
				Line:   1,
				Column: 1,
				Content: []*Node{{
					Kind:        ScalarNode,
					Tag:         "!!str",
					Value:       "a",
					LineComment: "# AI",
					Line:        1,
					Column:      1,
				}, {
					Kind: SequenceNode,
					Tag:  "!!seq",
					Content: []*Node{{
						Kind:   ScalarNode,
						Tag:    "!!str",
						Value:  "b",
						Line:   2,
						Column: 5,
					}},
					Line:   2,
					Column: 3,
				}, {
					Kind:   ScalarNode,
					Tag:    "!!str",
					Value:  "c",
					Line:   3,
					Column: 1,
				}, {
					Kind: SequenceNode,
					Tag:  "!!seq",
					Content: []*Node{{
						Kind:   ScalarNode,
						Tag:    "!!str",
						Value:  "d",
						Line:   4,
						Column: 5,
					}},
					Line:   4,
					Column: 3,
				}},
			}},
		},
	}, {
		"[decode]a:\n  # HM\n  - # HB1\n    # HB2\n    b: # IB\n      c # IC\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   MappingNode,
				Tag:    "!!map",
				Line:   1,
				Column: 1,
				Content: []*Node{{
					Kind:   ScalarNode,
					Style:  0x0,
					Tag:    "!!str",
					Value:  "a",
					Line:   1,
					Column: 1,
				}, {
					Kind:   SequenceNode,
					Tag:    "!!seq",
					Line:   3,
					Column: 3,
					Content: []*Node{{
						Kind:        MappingNode,
						Tag:         "!!map",
						HeadComment: "# HM",
						Line:        5,
						Column:      5,
						Content: []*Node{{
							Kind:        ScalarNode,
							Tag:         "!!str",
							Value:       "b",
							HeadComment: "# HB1\n# HB2",
							LineComment: "# IB",
							Line:        5,
							Column:      5,
						}, {
							Kind:        ScalarNode,
							Tag:         "!!str",
							Value:       "c",
							LineComment: "# IC",
							Line:        6,
							Column:      7,
						}},
					}},
				}},
			}},
		},
	}, {
		// When encoding the value above, it loses b's inline comment.
		"[encode]a:\n  # HM\n  - # HB1\n    # HB2\n    b: c # IC\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   MappingNode,
				Tag:    "!!map",
				Line:   1,
				Column: 1,
				Content: []*Node{{
					Kind:   ScalarNode,
					Style:  0x0,
					Tag:    "!!str",
					Value:  "a",
					Line:   1,
					Column: 1,
				}, {
					Kind:   SequenceNode,
					Tag:    "!!seq",
					Line:   3,
					Column: 3,
					Content: []*Node{{
						Kind:        MappingNode,
						Tag:         "!!map",
						HeadComment: "# HM",
						Line:        5,
						Column:      5,
						Content: []*Node{{
							Kind:        ScalarNode,
							Tag:         "!!str",
							Value:       "b",
							HeadComment: "# HB1\n# HB2",
							LineComment: "# IB",
							Line:        5,
							Column:      5,
						}, {
							Kind:        ScalarNode,
							Tag:         "!!str",
							Value:       "c",
							LineComment: "# IC",
							Line:        6,
							Column:      7,
						}},
					}},
				}},
			}},
		},
	}, {
		// Multiple cases of comment inlining next to mapping keys.
		"a: | # IA\n  str\nb: >- # IB\n  str\nc: # IC\n  - str\nd: # ID\n  str:\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   MappingNode,
				Tag:    "!!map",
				Line:   1,
				Column: 1,
				Content: []*Node{{
					Kind:   ScalarNode,
					Tag:    "!!str",
					Value:  "a",
					Line:   1,
					Column: 1,
				}, {
					Kind:        ScalarNode,
					Style:       LiteralStyle,
					Tag:         "!!str",
					Value:       "str\n",
					LineComment: "# IA",
					Line:        1,
					Column:      4,
				}, {
					Kind:   ScalarNode,
					Tag:    "!!str",
					Value:  "b",
					Line:   3,
					Column: 1,
				}, {
					Kind:        ScalarNode,
					Style:       FoldedStyle,
					Tag:         "!!str",
					Value:       "str",
					LineComment: "# IB",
					Line:        3,
					Column:      4,
				}, {
					Kind:        ScalarNode,
					Tag:         "!!str",
					Value:       "c",
					LineComment: "# IC",
					Line:        5,
					Column:      1,
				}, {
					Kind:   SequenceNode,
					Tag:    "!!seq",
					Line:   6,
					Column: 3,
					Content: []*Node{{
						Kind:   ScalarNode,
						Tag:    "!!str",
						Value:  "str",
						Line:   6,
						Column: 5,
					}},
				}, {
					Kind:        ScalarNode,
					Tag:         "!!str",
					Value:       "d",
					LineComment: "# ID",
					Line:        7,
					Column:      1,
				}, {
					Kind:   MappingNode,
					Tag:    "!!map",
					Line:   8,
					Column: 3,
					Content: []*Node{{
						Kind:   ScalarNode,
						Tag:    "!!str",
						Value:  "str",
						Line:   8,
						Column: 3,
					}, {
						Kind:   ScalarNode,
						Tag:    "!!null",
						Line:   8,
						Column: 7,
					}},
				}},
			}},
		},
	}, {
		// Indentless sequence.
		"[decode]a:\n# HM\n- # HB1\n  # HB2\n  b: # IB\n    c # IC\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   MappingNode,
				Tag:    "!!map",
				Line:   1,
				Column: 1,
				Content: []*Node{{
					Kind:   ScalarNode,
					Tag:    "!!str",
					Value:  "a",
					Line:   1,
					Column: 1,
				}, {
					Kind:   SequenceNode,
					Tag:    "!!seq",
					Line:   3,
					Column: 1,
					Content: []*Node{{
						Kind:        MappingNode,
						Tag:         "!!map",
						HeadComment: "# HM",
						Line:        5,
						Column:      3,
						Content: []*Node{{
							Kind:        ScalarNode,
							Tag:         "!!str",
							Value:       "b",
							HeadComment: "# HB1\n# HB2",
							LineComment: "# IB",
							Line:        5,
							Column:      3,
						}, {
							Kind:        ScalarNode,
							Tag:         "!!str",
							Value:       "c",
							LineComment: "# IC",
							Line:        6,
							Column:      5,
						}},
					}},
				}},
			}},
		},
	}, {
		"- a\n- b\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   SequenceNode,
				Value:  "",
				Tag:    "!!seq",
				Line:   1,
				Column: 1,
				Content: []*Node{{
					Kind:   ScalarNode,
					Value:  "a",
					Tag:    "!!str",
					Line:   1,
					Column: 3,
				}, {
					Kind:   ScalarNode,
					Value:  "b",
					Tag:    "!!str",
					Line:   2,
					Column: 3,
				}},
			}},
		},
	}, {
		"- a\n- - b\n  - c\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   SequenceNode,
				Tag:    "!!seq",
				Line:   1,
				Column: 1,
				Content: []*Node{{
					Kind:   ScalarNode,
					Value:  "a",
					Tag:    "!!str",
					Line:   1,
					Column: 3,
				}, {
					Kind:   SequenceNode,
					Tag:    "!!seq",
					Line:   2,
					Column: 3,
					Content: []*Node{{
						Kind:   ScalarNode,
						Value:  "b",
						Tag:    "!!str",
						Line:   2,
						Column: 5,
					}, {
						Kind:   ScalarNode,
						Value:  "c",
						Tag:    "!!str",
						Line:   3,
						Column: 5,
					}},
				}},
			}},
		},
	}, {
		"[a, b]\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   SequenceNode,
				Style:  FlowStyle,
				Value:  "",
				Tag:    "!!seq",
				Line:   1,
				Column: 1,
				Content: []*Node{{
					Kind:   ScalarNode,
					Value:  "a",
					Tag:    "!!str",
					Line:   1,
					Column: 2,
				}, {
					Kind:   ScalarNode,
					Value:  "b",
					Tag:    "!!str",
					Line:   1,
					Column: 5,
				}},
			}},
		},
	}, {
		"- a\n- [b, c]\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   SequenceNode,
				Tag:    "!!seq",
				Line:   1,
				Column: 1,
				Content: []*Node{{
					Kind:   ScalarNode,
					Value:  "a",
					Tag:    "!!str",
					Line:   1,
					Column: 3,
				}, {
					Kind:   SequenceNode,
					Tag:    "!!seq",
					Style:  FlowStyle,
					Line:   2,
					Column: 3,
					Content: []*Node{{
						Kind:   ScalarNode,
						Value:  "b",
						Tag:    "!!str",
						Line:   2,
						Column: 4,
					}, {
						Kind:   ScalarNode,
						Value:  "c",
						Tag:    "!!str",
						Line:   2,
						Column: 7,
					}},
				}},
			}},
		},
	}, {
		"a: &x 1\nb: &y 2\nc: *x\nd: *y\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   MappingNode,
				Line:   1,
				Column: 1,
				Tag:    "!!map",
				Content: []*Node{{
					Kind:   ScalarNode,
					Value:  "a",
					Tag:    "!!str",
					Line:   1,
					Column: 1,
				},
					saveNode("x", &Node{
						Kind:   ScalarNode,
						Value:  "1",
						Tag:    "!!int",
						Anchor: "x",
						Line:   1,
						Column: 4,
					}),
					{
						Kind:   ScalarNode,
						Value:  "b",
						Tag:    "!!str",
						Line:   2,
						Column: 1,
					},
					saveNode("y", &Node{
						Kind:   ScalarNode,
						Value:  "2",
						Tag:    "!!int",
						Anchor: "y",
						Line:   2,
						Column: 4,
					}),
					{
						Kind:   ScalarNode,
						Value:  "c",
						Tag:    "!!str",
						Line:   3,
						Column: 1,
					}, {
						Kind:   AliasNode,
						Value:  "x",
						Alias:  dropNode("x"),
						Line:   3,
						Column: 4,
					}, {
						Kind:   ScalarNode,
						Value:  "d",
						Tag:    "!!str",
						Line:   4,
						Column: 1,
					}, {
						Kind:   AliasNode,
						Value:  "y",
						Tag:    "",
						Alias:  dropNode("y"),
						Line:   4,
						Column: 4,
					}},
			}},
		},
	}, {

		"# One\n# Two\ntrue # Three\n# Four\n# Five\n",
		Node{
			Kind:   DocumentNode,
			Line:   3,
			Column: 1,
			Content: []*Node{{
				Kind:        ScalarNode,
				Value:       "true",
				Tag:         "!!bool",
				Line:        3,
				Column:      1,
				HeadComment: "# One\n# Two",
				LineComment: "# Three",
				FootComment: "# Four\n# Five",
			}},
		},
	}, {

		"# š\ntrue # š\n",
		Node{
			Kind:   DocumentNode,
			Line:   2,
			Column: 1,
			Content: []*Node{{
				Kind:        ScalarNode,
				Value:       "true",
				Tag:         "!!bool",
				Line:        2,
				Column:      1,
				HeadComment: "# š",
				LineComment: "# š",
			}},
		},
	}, {

		"[decode]\n# One\n\n# Two\n\n# Three\ntrue # Four\n# Five\n\n# Six\n\n# Seven\n",
		Node{
			Kind:        DocumentNode,
			Line:        7,
			Column:      1,
			HeadComment: "# One\n\n# Two",
			FootComment: "# Six\n\n# Seven",
			Content: []*Node{{
				Kind:        ScalarNode,
				Value:       "true",
				Tag:         "!!bool",
				Line:        7,
				Column:      1,
				HeadComment: "# Three",
				LineComment: "# Four",
				FootComment: "# Five",
			}},
		},
	}, {
		// Write out the pound character if missing from comments.
		"[encode]# One\n# Two\ntrue # Three\n# Four\n# Five\n",
		Node{
			Kind:   DocumentNode,
			Line:   3,
			Column: 1,
			Content: []*Node{{
				Kind:        ScalarNode,
				Value:       "true",
				Tag:         "!!bool",
				Line:        3,
				Column:      1,
				HeadComment: "One\nTwo\n",
				LineComment: "Three\n",
				FootComment: "Four\nFive\n",
			}},
		},
	}, {
		"[encode]#   One\n#   Two\ntrue #   Three\n#   Four\n#   Five\n",
		Node{
			Kind:   DocumentNode,
			Line:   3,
			Column: 1,
			Content: []*Node{{
				Kind:        ScalarNode,
				Value:       "true",
				Tag:         "!!bool",
				Line:        3,
				Column:      1,
				HeadComment: "  One\n  Two",
				LineComment: "  Three",
				FootComment: "  Four\n  Five",
			}},
		},
	}, {
		"# DH1\n\n# DH2\n\n# H1\n# H2\ntrue # I\n# F1\n# F2\n\n# DF1\n\n# DF2\n",
		Node{
			Kind:        DocumentNode,
			Line:        7,
			Column:      1,
			HeadComment: "# DH1\n\n# DH2",
			FootComment: "# DF1\n\n# DF2",
			Content: []*Node{{
				Kind:        ScalarNode,
				Value:       "true",
				Tag:         "!!bool",
				Line:        7,
				Column:      1,
				HeadComment: "# H1\n# H2",
				LineComment: "# I",
				FootComment: "# F1\n# F2",
			}},
		},
	}, {
		"# DH1\n\n# DH2\n\n# HA1\n# HA2\nka: va # IA\n# FA1\n# FA2\n\n# HB1\n# HB2\nkb: vb # IB\n# FB1\n# FB2\n\n# DF1\n\n# DF2\n",
		Node{
			Kind:        DocumentNode,
			Line:        7,
			Column:      1,
			HeadComment: "# DH1\n\n# DH2",
			FootComment: "# DF1\n\n# DF2",
			Content: []*Node{{
				Kind:   MappingNode,
				Tag:    "!!map",
				Line:   7,
				Column: 1,
				Content: []*Node{{
					Kind:        ScalarNode,
					Line:        7,
					Column:      1,
					Tag:         "!!str",
					Value:       "ka",
					HeadComment: "# HA1\n# HA2",
					FootComment: "# FA1\n# FA2",
				}, {
					Kind:        ScalarNode,
					Line:        7,
					Column:      5,
					Tag:         "!!str",
					Value:       "va",
					LineComment: "# IA",
				}, {
					Kind:        ScalarNode,
					Line:        13,
					Column:      1,
					Tag:         "!!str",
					Value:       "kb",
					HeadComment: "# HB1\n# HB2",
					FootComment: "# FB1\n# FB2",
				}, {
					Kind:        ScalarNode,
					Line:        13,
					Column:      5,
					Tag:         "!!str",
					Value:       "vb",
					LineComment: "# IB",
				}},
			}},
		},
	}, {
		"# DH1\n\n# DH2\n\n# HA1\n# HA2\n- la # IA\n# FA1\n# FA2\n\n# HB1\n# HB2\n- lb # IB\n# FB1\n# FB2\n\n# DF1\n\n# DF2\n",
		Node{
			Kind:        DocumentNode,
			Line:        7,
			Column:      1,
			HeadComment: "# DH1\n\n# DH2",
			FootComment: "# DF1\n\n# DF2",
			Content: []*Node{{
				Kind:   SequenceNode,
				Tag:    "!!seq",
				Line:   7,
				Column: 1,
				Content: []*Node{{
					Kind:        ScalarNode,
					Tag:         "!!str",
					Line:        7,
					Column:      3,
					Value:       "la",
					HeadComment: "# HA1\n# HA2",
					LineComment: "# IA",
					FootComment: "# FA1\n# FA2",
				}, {
					Kind:        ScalarNode,
					Tag:         "!!str",
					Line:        13,
					Column:      3,
					Value:       "lb",
					HeadComment: "# HB1\n# HB2",
					LineComment: "# IB",
					FootComment: "# FB1\n# FB2",
				}},
			}},
		},
	}, {
		"# DH1\n\n- la # IA\n# HB1\n- lb\n",
		Node{
			Kind:        DocumentNode,
			Line:        3,
			Column:      1,
			HeadComment: "# DH1",
			Content: []*Node{{
				Kind:   SequenceNode,
				Tag:    "!!seq",
				Line:   3,
				Column: 1,
				Content: []*Node{{
					Kind:        ScalarNode,
					Tag:         "!!str",
					Line:        3,
					Column:      3,
					Value:       "la",
					LineComment: "# IA",
				}, {
					Kind:        ScalarNode,
					Tag:         "!!str",
					Line:        5,
					Column:      3,
					Value:       "lb",
					HeadComment: "# HB1",
				}},
			}},
		},
	}, {
		"- la # IA\n- lb # IB\n- lc # IC\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   SequenceNode,
				Tag:    "!!seq",
				Line:   1,
				Column: 1,
				Content: []*Node{{
					Kind:        ScalarNode,
					Tag:         "!!str",
					Line:        1,
					Column:      3,
					Value:       "la",
					LineComment: "# IA",
				}, {
					Kind:        ScalarNode,
					Tag:         "!!str",
					Line:        2,
					Column:      3,
					Value:       "lb",
					LineComment: "# IB",
				}, {
					Kind:        ScalarNode,
					Tag:         "!!str",
					Line:        3,
					Column:      3,
					Value:       "lc",
					LineComment: "# IC",
				}},
			}},
		},
	}, {
		"# DH1\n\n# HL1\n- - la\n  # HB1\n  - lb\n",
		Node{
			Kind:        DocumentNode,
			Line:        4,
			Column:      1,
			HeadComment: "# DH1",
			Content: []*Node{{
				Kind:   SequenceNode,
				Tag:    "!!seq",
				Line:   4,
				Column: 1,
				Content: []*Node{{
					Kind:        SequenceNode,
					Tag:         "!!seq",
					Line:        4,
					Column:      3,
					HeadComment: "# HL1",
					Content: []*Node{{
						Kind:   ScalarNode,
						Tag:    "!!str",
						Line:   4,
						Column: 5,
						Value:  "la",
					}, {
						Kind:        ScalarNode,
						Tag:         "!!str",
						Line:        6,
						Column:      5,
						Value:       "lb",
						HeadComment: "# HB1",
					}},
				}},
			}},
		},
	}, {
		"# DH1\n\n# HL1\n- # HA1\n  - la\n  # HB1\n  - lb\n",
		Node{
			Kind:        DocumentNode,
			Line:        4,
			Column:      1,
			HeadComment: "# DH1",
			Content: []*Node{{
				Kind:   SequenceNode,
				Tag:    "!!seq",
				Line:   4,
				Column: 1,
				Content: []*Node{{
					Kind:        SequenceNode,
					Tag:         "!!seq",
					Line:        5,
					Column:      3,
					HeadComment: "# HL1",
					Content: []*Node{{
						Kind:        ScalarNode,
						Tag:         "!!str",
						Line:        5,
						Column:      5,
						Value:       "la",
						HeadComment: "# HA1",
					}, {
						Kind:        ScalarNode,
						Tag:         "!!str",
						Line:        7,
						Column:      5,
						Value:       "lb",
						HeadComment: "# HB1",
					}},
				}},
			}},
		},
	}, {
		"[decode]# DH1\n\n# HL1\n- # HA1\n\n  - la\n  # HB1\n  - lb\n",
		Node{
			Kind:        DocumentNode,
			Line:        4,
			Column:      1,
			HeadComment: "# DH1",
			Content: []*Node{{
				Kind:   SequenceNode,
				Tag:    "!!seq",
				Line:   4,
				Column: 1,
				Content: []*Node{{
					Kind:        SequenceNode,
					Tag:         "!!seq",
					Line:        6,
					Column:      3,
					HeadComment: "# HL1",
					Content: []*Node{{
						Kind:        ScalarNode,
						Tag:         "!!str",
						Line:        6,
						Column:      5,
						Value:       "la",
						HeadComment: "# HA1\n",
					}, {
						Kind:        ScalarNode,
						Tag:         "!!str",
						Line:        8,
						Column:      5,
						Value:       "lb",
						HeadComment: "# HB1",
					}},
				}},
			}},
		},
	}, {
		"# DH1\n\n# HA1\nka:\n  # HB1\n  kb:\n    # HC1\n    # HC2\n    - lc # IC\n    # FC1\n    # FC2\n\n    # HD1\n    - ld # ID\n    # FD1\n\n# DF1\n",
		Node{
			Kind:        DocumentNode,
			Line:        4,
			Column:      1,
			HeadComment: "# DH1",
			FootComment: "# DF1",
			Content: []*Node{{
				Kind:   MappingNode,
				Tag:    "!!map",
				Line:   4,
				Column: 1,
				Content: []*Node{{
					Kind:        ScalarNode,
					Tag:         "!!str",
					Line:        4,
					Column:      1,
					Value:       "ka",
					HeadComment: "# HA1",
				}, {
					Kind:   MappingNode,
					Tag:    "!!map",
					Line:   6,
					Column: 3,
					Content: []*Node{{
						Kind:        ScalarNode,
						Tag:         "!!str",
						Line:        6,
						Column:      3,
						Value:       "kb",
						HeadComment: "# HB1",
					}, {
						Kind:   SequenceNode,
						Line:   9,
						Column: 5,
						Tag:    "!!seq",
						Content: []*Node{{
							Kind:        ScalarNode,
							Tag:         "!!str",
							Line:        9,
							Column:      7,
							Value:       "lc",
							HeadComment: "# HC1\n# HC2",
							LineComment: "# IC",
							FootComment: "# FC1\n# FC2",
						}, {
							Kind:        ScalarNode,
							Tag:         "!!str",
							Line:        14,
							Column:      7,
							Value:       "ld",
							HeadComment: "# HD1",

							LineComment: "# ID",
							FootComment: "# FD1",
						}},
					}},
				}},
			}},
		},
	}, {
		"# DH1\n\n# HA1\nka:\n  # HB1\n  kb:\n    # HC1\n    # HC2\n    - lc # IC\n    # FC1\n    # FC2\n\n    # HD1\n    - ld # ID\n    # FD1\nke: ve\n\n# DF1\n",
		Node{
			Kind:        DocumentNode,
			Line:        4,
			Column:      1,
			HeadComment: "# DH1",
			FootComment: "# DF1",
			Content: []*Node{{
				Kind:   MappingNode,
				Tag:    "!!map",
				Line:   4,
				Column: 1,
				Content: []*Node{{
					Kind:        ScalarNode,
					Tag:         "!!str",
					Line:        4,
					Column:      1,
					Value:       "ka",
					HeadComment: "# HA1",
				}, {
					Kind:   MappingNode,
					Tag:    "!!map",
					Line:   6,
					Column: 3,
					Content: []*Node{{
						Kind:        ScalarNode,
						Tag:         "!!str",
						Line:        6,
						Column:      3,
						Value:       "kb",
						HeadComment: "# HB1",
					}, {
						Kind:   SequenceNode,
						Line:   9,
						Column: 5,
						Tag:    "!!seq",
						Content: []*Node{{
							Kind:        ScalarNode,
							Tag:         "!!str",
							Line:        9,
							Column:      7,
							Value:       "lc",
							HeadComment: "# HC1\n# HC2",
							LineComment: "# IC",
							FootComment: "# FC1\n# FC2",
						}, {
							Kind:        ScalarNode,
							Tag:         "!!str",
							Line:        14,
							Column:      7,
							Value:       "ld",
							HeadComment: "# HD1",
							LineComment: "# ID",
							FootComment: "# FD1",
						}},
					}},
				}, {
					Kind:   ScalarNode,
					Tag:    "!!str",
					Line:   16,
					Column: 1,
					Value:  "ke",
				}, {
					Kind:   ScalarNode,
					Tag:    "!!str",
					Line:   16,
					Column: 5,
					Value:  "ve",
				}},
			}},
		},
	}, {
		"# DH1\n\n# DH2\n\n# HA1\n# HA2\nka:\n  # HB1\n  # HB2\n  kb:\n" +
			"    # HC1\n    # HC2\n    kc:\n      # HD1\n      # HD2\n      kd: vd\n      # FD1\n      # FD2\n" +
			"    # FC1\n    # FC2\n  # FB1\n  # FB2\n# FA1\n# FA2\n\n# HE1\n# HE2\nke: ve\n# FE1\n# FE2\n\n# DF1\n\n# DF2\n",
		Node{
			Kind:        DocumentNode,
			HeadComment: "# DH1\n\n# DH2",
			FootComment: "# DF1\n\n# DF2",
			Line:        7,
			Column:      1,
			Content: []*Node{{
				Kind:   MappingNode,
				Tag:    "!!map",
				Line:   7,
				Column: 1,
				Content: []*Node{{
					Kind:        ScalarNode,
					Tag:         "!!str",
					Value:       "ka",
					HeadComment: "# HA1\n# HA2",
					FootComment: "# FA1\n# FA2",
					Line:        7,
					Column:      1,
				}, {
					Kind:   MappingNode,
					Tag:    "!!map",
					Line:   10,
					Column: 3,
					Content: []*Node{{
						Kind:        ScalarNode,
						Tag:         "!!str",
						Value:       "kb",
						HeadComment: "# HB1\n# HB2",
						FootComment: "# FB1\n# FB2",
						Line:        10,
						Column:      3,
					}, {
						Kind:   MappingNode,
						Tag:    "!!map",
						Line:   13,
						Column: 5,
						Content: []*Node{{
							Kind:        ScalarNode,
							Tag:         "!!str",
							Value:       "kc",
							HeadComment: "# HC1\n# HC2",
							FootComment: "# FC1\n# FC2",
							Line:        13,
							Column:      5,
						}, {
							Kind:   MappingNode,
							Tag:    "!!map",
							Line:   16,
							Column: 7,
							Content: []*Node{{
								Kind:        ScalarNode,
								Tag:         "!!str",
								Value:       "kd",
								HeadComment: "# HD1\n# HD2",
								FootComment: "# FD1\n# FD2",
								Line:        16,
								Column:      7,
							}, {
								Kind:   ScalarNode,
								Tag:    "!!str",
								Value:  "vd",
								Line:   16,
								Column: 11,
							}},
						}},
					}},
				}, {
					Kind:        ScalarNode,
					Tag:         "!!str",
					Value:       "ke",
					HeadComment: "# HE1\n# HE2",
					FootComment: "# FE1\n# FE2",
					Line:        28,
					Column:      1,
				}, {
					Kind:   ScalarNode,
					Tag:    "!!str",
					Value:  "ve",
					Line:   28,
					Column: 5,
				}},
			}},
		},
	}, {
		// Same as above but indenting ke in so it's also part of ka's value.
		"# DH1\n\n# DH2\n\n# HA1\n# HA2\nka:\n  # HB1\n  # HB2\n  kb:\n" +
			"    # HC1\n    # HC2\n    kc:\n      # HD1\n      # HD2\n      kd: vd\n      # FD1\n      # FD2\n" +
			"    # FC1\n    # FC2\n  # FB1\n  # FB2\n\n  # HE1\n  # HE2\n  ke: ve\n  # FE1\n  # FE2\n# FA1\n# FA2\n\n# DF1\n\n# DF2\n",
		Node{
			Kind:        DocumentNode,
			HeadComment: "# DH1\n\n# DH2",
			FootComment: "# DF1\n\n# DF2",
			Line:        7,
			Column:      1,
			Content: []*Node{{
				Kind:   MappingNode,
				Tag:    "!!map",
				Line:   7,
				Column: 1,
				Content: []*Node{{
					Kind:        ScalarNode,
					Tag:         "!!str",
					Value:       "ka",
					HeadComment: "# HA1\n# HA2",
					FootComment: "# FA1\n# FA2",
					Line:        7,
					Column:      1,
				}, {
					Kind:   MappingNode,
					Tag:    "!!map",
					Line:   10,
					Column: 3,
					Content: []*Node{{
						Kind:        ScalarNode,
						Tag:         "!!str",
						Value:       "kb",
						HeadComment: "# HB1\n# HB2",
						FootComment: "# FB1\n# FB2",
						Line:        10,
						Column:      3,
					}, {
						Kind:   MappingNode,
						Tag:    "!!map",
						Line:   13,
						Column: 5,
						Content: []*Node{{
							Kind:        ScalarNode,
							Tag:         "!!str",
							Value:       "kc",
							HeadComment: "# HC1\n# HC2",
							FootComment: "# FC1\n# FC2",
							Line:        13,
							Column:      5,
						}, {
							Kind:   MappingNode,
							Tag:    "!!map",
							Line:   16,
							Column: 7,
							Content: []*Node{{
								Kind:        ScalarNode,
								Tag:         "!!str",
								Value:       "kd",
								HeadComment: "# HD1\n# HD2",
								FootComment: "# FD1\n# FD2",
								Line:        16,
								Column:      7,
							}, {
								Kind:   ScalarNode,
								Tag:    "!!str",
								Value:  "vd",
								Line:   16,
								Column: 11,
							}},
						}},
					}, {
						Kind:        ScalarNode,
						Tag:         "!!str",
						Value:       "ke",
						HeadComment: "# HE1\n# HE2",
						FootComment: "# FE1\n# FE2",
						Line:        26,
						Column:      3,
					}, {
						Kind:   ScalarNode,
						Tag:    "!!str",
						Value:  "ve",
						Line:   26,
						Column: 7,
					}},
				}},
			}},
		},
	}, {
		// Decode only due to lack of newline at the end.
		"[decode]# HA1\nka:\n  # HB1\n  kb: vb\n  # FB1\n# FA1",
		Node{
			Kind:   DocumentNode,
			Line:   2,
			Column: 1,
			Content: []*Node{{
				Kind:   MappingNode,
				Tag:    "!!map",
				Line:   2,
				Column: 1,
				Content: []*Node{{
					Kind:        ScalarNode,
					Tag:         "!!str",
					Value:       "ka",
					HeadComment: "# HA1",
					FootComment: "# FA1",
					Line:        2,
					Column:      1,
				}, {
					Kind:   MappingNode,
					Tag:    "!!map",
					Line:   4,
					Column: 3,
					Content: []*Node{{
						Kind:        ScalarNode,
						Tag:         "!!str",
						Value:       "kb",
						HeadComment: "# HB1",
						FootComment: "# FB1",
						Line:        4,
						Column:      3,
					}, {
						Kind:   ScalarNode,
						Tag:    "!!str",
						Value:  "vb",
						Line:   4,
						Column: 7,
					}},
				}},
			}},
		},
	}, {
		// Same as above, but with newline at the end.
		"# HA1\nka:\n  # HB1\n  kb: vb\n  # FB1\n# FA1\n",
		Node{
			Kind:   DocumentNode,
			Line:   2,
			Column: 1,
			Content: []*Node{{
				Kind:   MappingNode,
				Tag:    "!!map",
				Line:   2,
				Column: 1,
				Content: []*Node{{
					Kind:        ScalarNode,
					Tag:         "!!str",
					Value:       "ka",
					HeadComment: "# HA1",
					FootComment: "# FA1",
					Line:        2,
					Column:      1,
				}, {
					Kind:   MappingNode,
					Tag:    "!!map",
					Line:   4,
					Column: 3,
					Content: []*Node{{
						Kind:        ScalarNode,
						Tag:         "!!str",
						Value:       "kb",
						HeadComment: "# HB1",
						FootComment: "# FB1",
						Line:        4,
						Column:      3,
					}, {
						Kind:   ScalarNode,
						Tag:    "!!str",
						Value:  "vb",
						Line:   4,
						Column: 7,
					}},
				}},
			}},
		},
	}, {
		// Same as above, but without FB1.
		"# HA1\nka:\n  # HB1\n  kb: vb\n# FA1\n",
		Node{
			Kind:   DocumentNode,
			Line:   2,
			Column: 1,
			Content: []*Node{{
				Kind:   MappingNode,
				Tag:    "!!map",
				Line:   2,
				Column: 1,
				Content: []*Node{{
					Kind:        ScalarNode,
					Tag:         "!!str",
					Value:       "ka",
					HeadComment: "# HA1",
					FootComment: "# FA1",
					Line:        2,
					Column:      1,
				}, {
					Kind:   MappingNode,
					Tag:    "!!map",
					Line:   4,
					Column: 3,
					Content: []*Node{{
						Kind:        ScalarNode,
						Tag:         "!!str",
						Value:       "kb",
						HeadComment: "# HB1",
						Line:        4,
						Column:      3,
					}, {
						Kind:   ScalarNode,
						Tag:    "!!str",
						Value:  "vb",
						Line:   4,
						Column: 7,
					}},
				}},
			}},
		},
	}, {
		// Same as above, but with two newlines at the end. Decode-only for that.
		"[decode]# HA1\nka:\n  # HB1\n  kb: vb\n  # FB1\n# FA1\n\n",
		Node{
			Kind:   DocumentNode,
			Line:   2,
			Column: 1,
			Content: []*Node{{
				Kind:   MappingNode,
				Tag:    "!!map",
				Line:   2,
				Column: 1,
				Content: []*Node{{
					Kind:        ScalarNode,
					Tag:         "!!str",
					Value:       "ka",
					HeadComment: "# HA1",
					FootComment: "# FA1",
					Line:        2,
					Column:      1,
				}, {
					Kind:   MappingNode,
					Tag:    "!!map",
					Line:   4,
					Column: 3,
					Content: []*Node{{
						Kind:        ScalarNode,
						Tag:         "!!str",
						Value:       "kb",
						HeadComment: "# HB1",
						FootComment: "# FB1",
						Line:        4,
						Column:      3,
					}, {
						Kind:   ScalarNode,
						Tag:    "!!str",
						Value:  "vb",
						Line:   4,
						Column: 7,
					}},
				}},
			}},
		},
	}, {
		// Similar to above, but make HB1 look more like a footer of ka.
		"[decode]# HA1\nka:\n# HB1\n\n  kb: vb\n# FA1\n",
		Node{
			Kind:   DocumentNode,
			Line:   2,
			Column: 1,
			Content: []*Node{{
				Kind:   MappingNode,
				Tag:    "!!map",
				Line:   2,
				Column: 1,
				Content: []*Node{{
					Kind:        ScalarNode,
					Tag:         "!!str",
					Value:       "ka",
					HeadComment: "# HA1",
					FootComment: "# FA1",
					Line:        2,
					Column:      1,
				}, {
					Kind:   MappingNode,
					Tag:    "!!map",
					Line:   5,
					Column: 3,
					Content: []*Node{{
						Kind:        ScalarNode,
						Tag:         "!!str",
						Value:       "kb",
						HeadComment: "# HB1\n",
						Line:        5,
						Column:      3,
					}, {
						Kind:   ScalarNode,
						Tag:    "!!str",
						Value:  "vb",
						Line:   5,
						Column: 7,
					}},
				}},
			}},
		},
	}, {
		"ka:\n  kb: vb\n# FA1\n\nkc: vc\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   MappingNode,
				Tag:    "!!map",
				Line:   1,
				Column: 1,
				Content: []*Node{{
					Kind:        ScalarNode,
					Tag:         "!!str",
					Value:       "ka",
					Line:        1,
					Column:      1,
					FootComment: "# FA1",
				}, {
					Kind:   MappingNode,
					Tag:    "!!map",
					Line:   2,
					Column: 3,
					Content: []*Node{{
						Kind:   ScalarNode,
						Tag:    "!!str",
						Value:  "kb",
						Line:   2,
						Column: 3,
					}, {
						Kind:   ScalarNode,
						Tag:    "!!str",
						Value:  "vb",
						Line:   2,
						Column: 7,
					}},
				}, {
					Kind:   ScalarNode,
					Tag:    "!!str",
					Value:  "kc",
					Line:   5,
					Column: 1,
				}, {
					Kind:   ScalarNode,
					Tag:    "!!str",
					Value:  "vc",
					Line:   5,
					Column: 5,
				}},
			}},
		},
	}, {
		"ka:\n  kb: vb\n# HC1\nkc: vc\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   MappingNode,
				Tag:    "!!map",
				Line:   1,
				Column: 1,
				Content: []*Node{{
					Kind:   ScalarNode,
					Tag:    "!!str",
					Value:  "ka",
					Line:   1,
					Column: 1,
				}, {
					Kind:   MappingNode,
					Tag:    "!!map",
					Line:   2,
					Column: 3,
					Content: []*Node{{
						Kind:   ScalarNode,
						Tag:    "!!str",
						Value:  "kb",
						Line:   2,
						Column: 3,
					}, {
						Kind:   ScalarNode,
						Tag:    "!!str",
						Value:  "vb",
						Line:   2,
						Column: 7,
					}},
				}, {
					Kind:        ScalarNode,
					Tag:         "!!str",
					Value:       "kc",
					HeadComment: "# HC1",
					Line:        4,
					Column:      1,
				}, {
					Kind:   ScalarNode,
					Tag:    "!!str",
					Value:  "vc",
					Line:   4,
					Column: 5,
				}},
			}},
		},
	}, {
		// Decode only due to empty line before HC1.
		"[decode]ka:\n  kb: vb\n\n# HC1\nkc: vc\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   MappingNode,
				Tag:    "!!map",
				Line:   1,
				Column: 1,
				Content: []*Node{{
					Kind:   ScalarNode,
					Tag:    "!!str",
					Value:  "ka",
					Line:   1,
					Column: 1,
				}, {
					Kind:   MappingNode,
					Tag:    "!!map",
					Line:   2,
					Column: 3,
					Content: []*Node{{
						Kind:   ScalarNode,
						Tag:    "!!str",
						Value:  "kb",
						Line:   2,
						Column: 3,
					}, {
						Kind:   ScalarNode,
						Tag:    "!!str",
						Value:  "vb",
						Line:   2,
						Column: 7,
					}},
				}, {
					Kind:        ScalarNode,
					Tag:         "!!str",
					Value:       "kc",
					HeadComment: "# HC1",
					Line:        5,
					Column:      1,
				}, {
					Kind:   ScalarNode,
					Tag:    "!!str",
					Value:  "vc",
					Line:   5,
					Column: 5,
				}},
			}},
		},
	}, {
		// Decode-only due to empty lines around HC1.
		"[decode]ka:\n  kb: vb\n\n# HC1\n\nkc: vc\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   MappingNode,
				Tag:    "!!map",
				Line:   1,
				Column: 1,
				Content: []*Node{{
					Kind:   ScalarNode,
					Tag:    "!!str",
					Value:  "ka",
					Line:   1,
					Column: 1,
				}, {
					Kind:   MappingNode,
					Tag:    "!!map",
					Line:   2,
					Column: 3,
					Content: []*Node{{
						Kind:   ScalarNode,
						Tag:    "!!str",
						Value:  "kb",
						Line:   2,
						Column: 3,
					}, {
						Kind:   ScalarNode,
						Tag:    "!!str",
						Value:  "vb",
						Line:   2,
						Column: 7,
					}},
				}, {
					Kind:        ScalarNode,
					Tag:         "!!str",
					Value:       "kc",
					HeadComment: "# HC1\n",
					Line:        6,
					Column:      1,
				}, {
					Kind:   ScalarNode,
					Tag:    "!!str",
					Value:  "vc",
					Line:   6,
					Column: 5,
				}},
			}},
		},
	}, {
		"ka: # IA\n  kb: # IB\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   MappingNode,
				Tag:    "!!map",
				Line:   1,
				Column: 1,
				Content: []*Node{{
					Kind:        ScalarNode,
					Tag:         "!!str",
					Value:       "ka",
					Line:        1,
					Column:      1,
					LineComment: "# IA",
				}, {
					Kind:   MappingNode,
					Tag:    "!!map",
					Line:   2,
					Column: 3,
					Content: []*Node{{
						Kind:        ScalarNode,
						Tag:         "!!str",
						Value:       "kb",
						Line:        2,
						Column:      3,
						LineComment: "# IB",
					}, {
						Kind:   ScalarNode,
						Tag:    "!!null",
						Line:   2,
						Column: 6,
					}},
				}},
			}},
		},
	}, {
		"# HA1\nka:\n  # HB1\n  kb: vb\n  # FB1\n# HC1\n# HC2\nkc: vc\n# FC1\n# FC2\n",
		Node{
			Kind:   DocumentNode,
			Line:   2,
			Column: 1,
			Content: []*Node{{
				Kind:   MappingNode,
				Tag:    "!!map",
				Line:   2,
				Column: 1,
				Content: []*Node{{
					Kind:        ScalarNode,
					Tag:         "!!str",
					Value:       "ka",
					HeadComment: "# HA1",
					Line:        2,
					Column:      1,
				}, {
					Kind:   MappingNode,
					Tag:    "!!map",
					Line:   4,
					Column: 3,
					Content: []*Node{{
						Kind:        ScalarNode,
						Tag:         "!!str",
						Value:       "kb",
						HeadComment: "# HB1",
						FootComment: "# FB1",
						Line:        4,
						Column:      3,
					}, {
						Kind:   ScalarNode,
						Tag:    "!!str",
						Value:  "vb",
						Line:   4,
						Column: 7,
					}},
				}, {
					Kind:        ScalarNode,
					Tag:         "!!str",
					Value:       "kc",
					HeadComment: "# HC1\n# HC2",
					FootComment: "# FC1\n# FC2",
					Line:        8,
					Column:      1,
				}, {
					Kind:   ScalarNode,
					Tag:    "!!str",
					Value:  "vc",
					Line:   8,
					Column: 5,
				}},
			}},
		},
	}, {
		// Same as above, but decode only due to empty line between ka's value and kc's headers.
		"[decode]# HA1\nka:\n  # HB1\n  kb: vb\n  # FB1\n\n# HC1\n# HC2\nkc: vc\n# FC1\n# FC2\n",
		Node{
			Kind:   DocumentNode,
			Line:   2,
			Column: 1,
			Content: []*Node{{
				Kind:   MappingNode,
				Tag:    "!!map",
				Line:   2,
				Column: 1,
				Content: []*Node{{
					Kind:        ScalarNode,
					Tag:         "!!str",
					Value:       "ka",
					HeadComment: "# HA1",
					Line:        2,
					Column:      1,
				}, {
					Kind:   MappingNode,
					Tag:    "!!map",
					Line:   4,
					Column: 3,
					Content: []*Node{{
						Kind:        ScalarNode,
						Tag:         "!!str",
						Value:       "kb",
						HeadComment: "# HB1",
						FootComment: "# FB1",
						Line:        4,
						Column:      3,
					}, {
						Kind:   ScalarNode,
						Tag:    "!!str",
						Value:  "vb",
						Line:   4,
						Column: 7,
					}},
				}, {
					Kind:        ScalarNode,
					Tag:         "!!str",
					Value:       "kc",
					HeadComment: "# HC1\n# HC2",
					FootComment: "# FC1\n# FC2",
					Line:        9,
					Column:      1,
				}, {
					Kind:   ScalarNode,
					Tag:    "!!str",
					Value:  "vc",
					Line:   9,
					Column: 5,
				}},
			}},
		},
	}, {
		"# H1\n[la, lb] # I\n# F1\n",
		Node{
			Kind:   DocumentNode,
			Line:   2,
			Column: 1,
			Content: []*Node{{
				Kind:        SequenceNode,
				Tag:         "!!seq",
				Style:       FlowStyle,
				Line:        2,
				Column:      1,
				HeadComment: "# H1",
				LineComment: "# I",
				FootComment: "# F1",
				Content: []*Node{{
					Kind:   ScalarNode,
					Tag:    "!!str",
					Line:   2,
					Column: 2,
					Value:  "la",
				}, {
					Kind:   ScalarNode,
					Tag:    "!!str",
					Line:   2,
					Column: 6,
					Value:  "lb",
				}},
			}},
		},
	}, {
		"# DH1\n\n# SH1\n[\n  # HA1\n  la, # IA\n  # FA1\n\n  # HB1\n  lb, # IB\n  # FB1\n]\n# SF1\n\n# DF1\n",
		Node{
			Kind:        DocumentNode,
			Line:        4,
			Column:      1,
			HeadComment: "# DH1",
			FootComment: "# DF1",
			Content: []*Node{{
				Kind:        SequenceNode,
				Tag:         "!!seq",
				Style:       FlowStyle,
				Line:        4,
				Column:      1,
				HeadComment: "# SH1",
				FootComment: "# SF1",
				Content: []*Node{{
					Kind:        ScalarNode,
					Tag:         "!!str",
					Line:        6,
					Column:      3,
					Value:       "la",
					HeadComment: "# HA1",
					LineComment: "# IA",
					FootComment: "# FA1",
				}, {
					Kind:        ScalarNode,
					Tag:         "!!str",
					Line:        10,
					Column:      3,
					Value:       "lb",
					HeadComment: "# HB1",
					LineComment: "# IB",
					FootComment: "# FB1",
				}},
			}},
		},
	}, {
		// Same as above, but with extra newlines before FB1 and FB2
		"[decode]# DH1\n\n# SH1\n[\n  # HA1\n  la, # IA\n  # FA1\n\n  # HB1\n  lb, # IB\n\n\n  # FB1\n\n# FB2\n]\n# SF1\n\n# DF1\n",
		Node{
			Kind:        DocumentNode,
			Line:        4,
			Column:      1,
			HeadComment: "# DH1",
			FootComment: "# DF1",
			Content: []*Node{{
				Kind:        SequenceNode,
				Tag:         "!!seq",
				Style:       FlowStyle,
				Line:        4,
				Column:      1,
				HeadComment: "# SH1",
				FootComment: "# SF1",
				Content: []*Node{{
					Kind:        ScalarNode,
					Tag:         "!!str",
					Line:        6,
					Column:      3,
					Value:       "la",
					HeadComment: "# HA1",
					LineComment: "# IA",
					FootComment: "# FA1",
				}, {
					Kind:        ScalarNode,
					Tag:         "!!str",
					Line:        10,
					Column:      3,
					Value:       "lb",
					HeadComment: "# HB1",
					LineComment: "# IB",
					FootComment: "# FB1\n\n# FB2",
				}},
			}},
		},
	}, {
		"# DH1\n\n# SH1\n[\n  # HA1\n  la,\n  # FA1\n\n  # HB1\n  lb,\n  # FB1\n]\n# SF1\n\n# DF1\n",
		Node{
			Kind:        DocumentNode,
			Line:        4,
			Column:      1,
			HeadComment: "# DH1",
			FootComment: "# DF1",
			Content: []*Node{{
				Kind:        SequenceNode,
				Tag:         "!!seq",
				Style:       FlowStyle,
				Line:        4,
				Column:      1,
				HeadComment: "# SH1",
				FootComment: "# SF1",
				Content: []*Node{{
					Kind:        ScalarNode,
					Tag:         "!!str",
					Line:        6,
					Column:      3,
					Value:       "la",
					HeadComment: "# HA1",
					FootComment: "# FA1",
				}, {
					Kind:        ScalarNode,
					Tag:         "!!str",
					Line:        10,
					Column:      3,
					Value:       "lb",
					HeadComment: "# HB1",
					FootComment: "# FB1",
				}},
			}},
		},
	}, {
		"ka:\n  kb: [\n    # HA1\n    la,\n    # FA1\n\n    # HB1\n    lb,\n    # FB1\n  ]\n",
		Node{
			Kind:   DocumentNode,
			Line:   1,
			Column: 1,
			Content: []*Node{{
				Kind:   MappingNode,
				Tag:    "!!map",
				Line:   1,
				Column: 1,
				Content: []*Node{{
					Kind:   ScalarNode,
					Tag:    "!!str",
					Value:  "ka",
					Line:   1,
					Column: 1,
				}, {
					Kind:   0x4,
					Tag:    "!!map",
					Line:   2,
					Column: 3,
					Content: []*Node{{
						Kind:   ScalarNode,
						Tag:    "!!str",
						Value:  "kb",
						Line:   2,
						Column: 3,
					}, {
						Kind:   SequenceNode,
						Style:  0x20,
						Tag:    "!!seq",
						Line:   2,
						Column: 7,
						Content: []*Node{{
							Kind:        ScalarNode,
							Tag:         "!!str",
							Value:       "la",
							HeadComment: "# HA1",
							FootComment: "# FA1",
							Line:        4,
							Column:      5,
						}, {
							Kind:        ScalarNode,
							Tag:         "!!str",
							Value:       "lb",
							HeadComment: "# HB1",
							FootComment: "# FB1",
							Line:        8,
							Column:      5,
						}},
					}},
				}},
			}},
		},
	}, {
		"# DH1\n\n# MH1\n{\n  # HA1\n  ka: va, # IA\n  # FA1\n\n  # HB1\n  kb: vb, # IB\n  # FB1\n}\n# MF1\n\n# DF1\n",
		Node{
			Kind:        DocumentNode,
			Line:        4,
			Column:      1,
			HeadComment: "# DH1",
			FootComment: "# DF1",
			Content: []*Node{{
				Kind:        MappingNode,
				Tag:         "!!map",
				Style:       FlowStyle,
				Line:        4,
				Column:      1,
				HeadComment: "# MH1",
				FootComment: "# MF1",
				Content: []*Node{{
					Kind:        ScalarNode,
					Tag:         "!!str",
					Line:        6,
					Column:      3,
					Value:       "ka",
					HeadComment: "# HA1",
					FootComment: "# FA1",
				}, {
					Kind:        ScalarNode,
					Tag:         "!!str",
					Line:        6,
					Column:      7,
					Value:       "va",
					LineComment: "# IA",
				}, {
					Kind:        ScalarNode,
					Tag:         "!!str",
					Line:        10,
					Column:      3,
					Value:       "kb",
					HeadComment: "# HB1",
					FootComment: "# FB1",
				}, {
					Kind:        ScalarNode,
					Tag:         "!!str",
					Line:        10,
					Column:      7,
					Value:       "vb",
					LineComment: "# IB",
				}},
			}},
		},
	}, {
		"# DH1\n\n# MH1\n{\n  # HA1\n  ka: va,\n  # FA1\n\n  # HB1\n  kb: vb,\n  # FB1\n}\n# MF1\n\n# DF1\n",
		Node{
			Kind:        DocumentNode,
			Line:        4,
			Column:      1,
			HeadComment: "# DH1",
			FootComment: "# DF1",
			Content: []*Node{{
				Kind:        MappingNode,
				Tag:         "!!map",
				Style:       FlowStyle,
				Line:        4,
				Column:      1,
				HeadComment: "# MH1",
				FootComment: "# MF1",
				Content: []*Node{{
					Kind:        ScalarNode,
					Tag:         "!!str",
					Line:        6,
					Column:      3,
					Value:       "ka",
					HeadComment: "# HA1",
					FootComment: "# FA1",
				}, {
					Kind:   ScalarNode,
					Tag:    "!!str",
					Line:   6,
					Column: 7,
					Value:  "va",
				}, {
					Kind:        ScalarNode,
					Tag:         "!!str",
					Line:        10,
					Column:      3,
					Value:       "kb",
					HeadComment: "# HB1",
					FootComment: "# FB1",
				}, {
					Kind:   ScalarNode,
					Tag:    "!!str",
					Line:   10,
					Column: 7,
					Value:  "vb",
				}},
			}},
		},
	}, {
		"# DH1\n\n# DH2\n\n# HA1\n# HA2\n- &x la # IA\n# FA1\n# FA2\n\n# HB1\n# HB2\n- *x # IB\n# FB1\n# FB2\n\n# DF1\n\n# DF2\n",
		Node{
			Kind:        DocumentNode,
			Line:        7,
			Column:      1,
			HeadComment: "# DH1\n\n# DH2",
			FootComment: "# DF1\n\n# DF2",
			Content: []*Node{{
				Kind:   SequenceNode,
				Tag:    "!!seq",
				Line:   7,
				Column: 1,
				Content: []*Node{
					saveNode("x", &Node{
						Kind:        ScalarNode,
						Tag:         "!!str",
						Line:        7,
						Column:      3,
						Value:       "la",
						HeadComment: "# HA1\n# HA2",
						LineComment: "# IA",
						FootComment: "# FA1\n# FA2",
						Anchor:      "x",
					}), {
						Kind:        AliasNode,
						Line:        13,
						Column:      3,
						Value:       "x",
						Alias:       dropNode("x"),
						HeadComment: "# HB1\n# HB2",
						LineComment: "# IB",
						FootComment: "# FB1\n# FB2",
					},
				},
			}},
		},
	},
}

func (s *S) TestNodeRoundtrip(c *C) {
	defer os.Setenv("TZ", os.Getenv("TZ"))
	os.Setenv("TZ", "UTC")
	for i, item := range nodeTests {
		c.Logf("test %d: %q", i, item.yaml)

		if strings.Contains(item.yaml, "#") {
			var buf bytes.Buffer
			fprintComments(&buf, &item.node, "    ")
			c.Logf("  expected comments:\n%s", buf.Bytes())
		}

		decode := true
		encode := true

		testYaml := item.yaml
		if s := strings.TrimPrefix(testYaml, "[decode]"); s != testYaml {
			encode = false
			testYaml = s
		}
		if s := strings.TrimPrefix(testYaml, "[encode]"); s != testYaml {
			decode = false
			testYaml = s
		}

		if decode {
			var node Node
			err := Unmarshal([]byte(testYaml), &node)
			c.Assert(err, IsNil)
			if strings.Contains(item.yaml, "#") {
				var buf bytes.Buffer
				fprintComments(&buf, &node, "    ")
				c.Logf("  obtained comments:\n%s", buf.Bytes())
			}
			c.Assert(&node, DeepEquals, &item.node)
		}
		if encode {
			node := deepCopyNode(&item.node, nil)
			buf := bytes.Buffer{}
			enc := NewEncoder(&buf)
			enc.SetIndent(2)
			err := enc.Encode(node)
			c.Assert(err, IsNil)
			err = enc.Close()
			c.Assert(err, IsNil)
			c.Assert(buf.String(), Equals, testYaml)

			// Ensure there were no mutations to the tree.
			c.Assert(node, DeepEquals, &item.node)
		}
	}
}

func deepCopyNode(node *Node, cache map[*Node]*Node) *Node {
	if n, ok := cache[node]; ok {
		return n
	}
	if cache == nil {
		cache = make(map[*Node]*Node)
	}
	copy := *node
	cache[node] = &copy
	copy.Content = nil
	for _, elem := range node.Content {
		copy.Content = append(copy.Content, deepCopyNode(elem, cache))
	}
	if node.Alias != nil {
		copy.Alias = deepCopyNode(node.Alias, cache)
	}
	return &copy
}

var savedNodes = make(map[string]*Node)

func saveNode(name string, node *Node) *Node {
	savedNodes[name] = node
	return node
}

func peekNode(name string) *Node {
	return savedNodes[name]
}

func dropNode(name string) *Node {
	node := savedNodes[name]
	delete(savedNodes, name)
	return node
}

var setStringTests = []struct {
	str  string
	yaml string
	node Node
}{
	{
		"something simple",
		"something simple\n",
		Node{
			Kind:  ScalarNode,
			Value: "something simple",
			Tag:   "!!str",
		},
	}, {
		`"quoted value"`,
		"'\"quoted value\"'\n",
		Node{
			Kind:  ScalarNode,
			Value: `"quoted value"`,
			Tag:   "!!str",
		},
	}, {
		"multi\nline",
		"|-\n  multi\n  line\n",
		Node{
			Kind:  ScalarNode,
			Value: "multi\nline",
			Tag:   "!!str",
			Style: LiteralStyle,
		},
	}, {
		"123",
		"\"123\"\n",
		Node{
			Kind:  ScalarNode,
			Value: "123",
			Tag:   "!!str",
		},
	}, {
		"multi\nline\n",
		"|\n  multi\n  line\n",
		Node{
			Kind:  ScalarNode,
			Value: "multi\nline\n",
			Tag:   "!!str",
			Style: LiteralStyle,
		},
	}, {
		"\x80\x81\x82",
		"!!binary gIGC\n",
		Node{
			Kind:  ScalarNode,
			Value: "gIGC",
			Tag:   "!!binary",
		},
	},
}

func (s *S) TestSetString(c *C) {
	defer os.Setenv("TZ", os.Getenv("TZ"))
	os.Setenv("TZ", "UTC")
	for i, item := range setStringTests {
		c.Logf("test %d: %q", i, item.str)

		var node Node

		node.SetString(item.str)

		c.Assert(node, DeepEquals, item.node)

		buf := bytes.Buffer{}
		enc := NewEncoder(&buf)
		enc.SetIndent(2)
		err := enc.Encode(&item.node)
		c.Assert(err, IsNil)
		err = enc.Close()
		c.Assert(err, IsNil)
		c.Assert(buf.String(), Equals, item.yaml)

		var doc Node
		err = Unmarshal([]byte(item.yaml), &doc)
		c.Assert(err, IsNil)

		var str string
		err = node.Decode(&str)
		c.Assert(err, IsNil)
		c.Assert(str, Equals, item.str)
	}
}

var nodeEncodeDecodeTests = []struct {
	value interface{}
	yaml  string
	node  Node
}{{
	"something simple",
	"something simple\n",
	Node{
		Kind:  ScalarNode,
		Value: "something simple",
		Tag:   "!!str",
	},
}, {
	`"quoted value"`,
	"'\"quoted value\"'\n",
	Node{
		Kind:  ScalarNode,
		Style: SingleQuotedStyle,
		Value: `"quoted value"`,
		Tag:   "!!str",
	},
}, {
	123,
	"123",
	Node{
		Kind:  ScalarNode,
		Value: `123`,
		Tag:   "!!int",
	},
}, {
	[]interface{}{1, 2},
	"[1, 2]",
	Node{
		Kind: SequenceNode,
		Tag:  "!!seq",
		Content: []*Node{{
			Kind:  ScalarNode,
			Value: "1",
			Tag:   "!!int",
		}, {
			Kind:  ScalarNode,
			Value: "2",
			Tag:   "!!int",
		}},
	},
}, {
	map[string]interface{}{"a": "b"},
	"a: b",
	Node{
		Kind: MappingNode,
		Tag:  "!!map",
		Content: []*Node{{
			Kind:  ScalarNode,
			Value: "a",
			Tag:   "!!str",
		}, {
			Kind:  ScalarNode,
			Value: "b",
			Tag:   "!!str",
		}},
	},
}}

func (s *S) TestNodeEncodeDecode(c *C) {
	for i, item := range nodeEncodeDecodeTests {
		c.Logf("Encode/Decode test value #%d: %#v", i, item.value)

		var v interface{}
		err := item.node.Decode(&v)
		c.Assert(err, IsNil)
		c.Assert(v, DeepEquals, item.value)

		var n Node
		err = n.Encode(item.value)
		c.Assert(err, IsNil)
		c.Assert(n, DeepEquals, item.node)
	}
}

func (s *S) TestNodeZeroEncodeDecode(c *C) {
	// Zero node value behaves as nil when encoding...
	var n Node
	data, err := Marshal(&n)
	c.Assert(err, IsNil)
	c.Assert(string(data), Equals, "null\n")

	// ... and decoding.
	var v *struct{} = &struct{}{}
	c.Assert(n.Decode(&v), IsNil)
	c.Assert(v, IsNil)

	// ... and even when looking for its tag.
	c.Assert(n.ShortTag(), Equals, "!!null")

	// Kind zero is still unknown, though.
	n.Line = 1
	_, err = Marshal(&n)
	c.Assert(err, ErrorMatches, "yaml: cannot encode node with unknown kind 0")
	c.Assert(n.Decode(&v), ErrorMatches, "yaml: cannot decode node with unknown kind 0")
}

func (s *S) TestNodeOmitEmpty(c *C) {
	var v struct {
		A int
		B Node ",omitempty"
	}
	v.A = 1
	data, err := Marshal(&v)
	c.Assert(err, IsNil)
	c.Assert(string(data), Equals, "a: 1\n")

	v.B.Line = 1
	_, err = Marshal(&v)
	c.Assert(err, ErrorMatches, "yaml: cannot encode node with unknown kind 0")
}

func fprintComments(out io.Writer, node *Node, indent string) {
	switch node.Kind {
	case ScalarNode:
		fmt.Fprintf(out, "%s<%s> ", indent, node.Value)
		fprintCommentSet(out, node)
		fmt.Fprintf(out, "\n")
	case DocumentNode:
		fmt.Fprintf(out, "%s<DOC> ", indent)
		fprintCommentSet(out, node)
		fmt.Fprintf(out, "\n")
		for i := 0; i < len(node.Content); i++ {
			fprintComments(out, node.Content[i], indent+"  ")
		}
	case MappingNode:
		fmt.Fprintf(out, "%s<MAP> ", indent)
		fprintCommentSet(out, node)
		fmt.Fprintf(out, "\n")
		for i := 0; i < len(node.Content); i += 2 {
			fprintComments(out, node.Content[i], indent+"  ")
			fprintComments(out, node.Content[i+1], indent+"  ")
		}
	case SequenceNode:
		fmt.Fprintf(out, "%s<SEQ> ", indent)
		fprintCommentSet(out, node)
		fmt.Fprintf(out, "\n")
		for i := 0; i < len(node.Content); i++ {
			fprintComments(out, node.Content[i], indent+"  ")
		}
	}
}

func fprintCommentSet(out io.Writer, node *Node) {
	if len(node.HeadComment)+len(node.LineComment)+len(node.FootComment) > 0 {
		fmt.Fprintf(out, "%q / %q / %q", node.HeadComment, node.LineComment, node.FootComment)
	}
}
