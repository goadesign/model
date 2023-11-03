package stz

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/kylelemons/godebug/diff"
)

func TestGet(t *testing.T) {
	const wID = "42"
	var wkspc = workspace(t)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/workspace/%s", wID) {
			t.Errorf("got path %s, expected %s", req.URL.String(), fmt.Sprintf("/workspace/%s", wID))
		}
		validateHeaders(t, req)
		rw.Header().Add("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(rw).Encode(wkspc); err != nil {
			t.Fatalf("failed to encode response: %s", err)
		}
	}))
	defer server.Close()

	// Substitute structurizr service host and scheme for tests.
	host := Host
	defer func() { Host = host }()
	u, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("failed to parse test server URL %q: %s", server.URL, err)
	}
	Host = u.Host
	scheme := Scheme
	defer func() { Scheme = scheme }()
	Scheme = "http"

	c := NewClient("key", "secret")
	wk, err := c.Get(wID)

	if err != nil {
		t.Errorf("Get failed with %s", err)
	}
	js, err := json.MarshalIndent(wk, "", "   ")
	if err != nil {
		t.Fatalf("failed to marshal response for comparison: %s", err)
	}
	js2, _ := json.MarshalIndent(wkspc, "", "   ")
	dif := diff.Diff(string(js), string(js2))
	if dif != "" {
		t.Errorf("invalid response content, got vs. expected:\n%s", dif)
	}
}

func TestAuth(t *testing.T) {
	wid, key, secret := config(t)
	c := NewClient(key, secret)
	_, err := c.Get(wid)
	if err != nil {
		t.Errorf("failed to retrieve workspace: %s", err)
	}
}

func validateHeaders(t *testing.T, req *http.Request) {
	if req.Header.Get("nonce") == "" {
		t.Errorf("missing nonce header")
	}
	if req.Header.Get("X-Authorization") == "" {
		t.Errorf("missing X-Authorization header")
	}
}

func workspace(t *testing.T) (workspace *Workspace) {
	err := json.Unmarshal([]byte(bigBankPLC), &workspace)
	if err != nil {
		t.Fatalf("unable to load test workspace: %s", err)
	}
	return
}

func config(t *testing.T) (workspaceID, key, secret string) {
	workspaceID = os.Getenv("STRUCTURIZR_WORKSPACE_ID")
	if workspaceID == "" {
		t.Skip("STRUCTURIZR_WORKSPACE_ID not set")
	}
	key = os.Getenv("STRUCTURIZR_KEY")
	if key == "" {
		t.Skip("STRUCTURIZR_KEY not set")
	}
	secret = os.Getenv("STRUCTURIZR_SECRET")
	if secret == "" {
		t.Skip("STRUCTURIZR_SECRET not set")
	}
	return
}

// Serialized workspace taken from
// https://raw.githubusercontent.com/structurizr/examples/main/json/big-bank-plc/workspace.json
// and replaced integer instances with string instances since [Structurizr DSL v1.22](https://github.com/structurizr/dsl/releases/tag/v1.22.0)
var bigBankPLC = `{
    "name": "Big Bank plc",
    "description": "This is an example workspace to illustrate the key features of Structurizr, via the DSL, based around a fictional online banking system.",
    "lastModifiedDate": "2022-02-28T16:16:49Z",
    "configuration": {},
    "model": {
        "enterprise": {
            "name": "Big Bank plc"
        },
        "people": [
            {
                "id": "3",
                "tags": "Element,Person,Bank Staff",
                "name": "Back Office Staff",
                "description": "Administration and support staff within the bank.",
                "relationships": [
                    {
                        "id": "27",
                        "tags": "Relationship",
                        "sourceId": "3",
                        "destinationId": "4",
                        "description": "Uses",
                        "defaultTags": [
                            "Relationship"
                        ]
                    }
                ],
                "location": "Internal",
                "defaultTags": [
                    "Element",
                    "Person"
                ]
            },
            {
                "id": "2",
                "tags": "Element,Person,Bank Staff",
                "name": "Customer Service Staff",
                "description": "Customer service staff within the bank.",
                "relationships": [
                    {
                        "id": "24",
                        "tags": "Relationship",
                        "sourceId": "2",
                        "destinationId": "4",
                        "description": "Uses",
                        "defaultTags": [
                            "Relationship"
                        ]
                    }
                ],
                "location": "Internal",
                "defaultTags": [
                    "Element",
                    "Person"
                ]
            },
            {
                "id": "1",
                "tags": "Element,Person,Customer",
                "name": "Personal Banking Customer",
                "description": "A customer of the bank, with personal bank accounts.",
                "relationships": [
                    {
                        "id": "19",
                        "tags": "Relationship",
                        "sourceId": "1",
                        "destinationId": "7",
                        "description": "Views account balances, and makes payments using",
                        "defaultTags": [
                            "Relationship"
                        ]
                    },
                    {
                        "id": "23",
                        "tags": "Relationship",
                        "sourceId": "1",
                        "destinationId": "2",
                        "description": "Asks questions to",
                        "technology": "Telephone",
                        "defaultTags": [
                            "Relationship"
                        ]
                    },
                    {
                        "id": "25",
                        "tags": "Relationship",
                        "sourceId": "1",
                        "destinationId": "6",
                        "description": "Withdraws cash using",
                        "defaultTags": [
                            "Relationship"
                        ]
                    },
                    {
                        "id": "28",
                        "tags": "Relationship",
                        "sourceId": "1",
                        "destinationId": "10",
                        "description": "Visits bigbank.com/ib using",
                        "technology": "HTTPS",
                        "defaultTags": [
                            "Relationship"
                        ]
                    },
                    {
                        "id": "29",
                        "tags": "Relationship",
                        "sourceId": "1",
                        "destinationId": "8",
                        "description": "Views account balances, and makes payments using",
                        "defaultTags": [
                            "Relationship"
                        ]
                    },
                    {
                        "id": "30",
                        "tags": "Relationship",
                        "sourceId": "1",
                        "destinationId": "9",
                        "description": "Views account balances, and makes payments using",
                        "defaultTags": [
                            "Relationship"
                        ]
                    }
                ],
                "location": "External",
                "defaultTags": [
                    "Element",
                    "Person"
                ]
            }
        ],
        "softwareSystems": [
            {
                "id": "6",
                "tags": "Element,Software System,Existing System",
                "name": "ATM",
                "description": "Allows customers to withdraw cash.",
                "relationships": [
                    {
                        "id": "26",
                        "tags": "Relationship",
                        "sourceId": "6",
                        "destinationId": "4",
                        "description": "Uses",
                        "defaultTags": [
                            "Relationship"
                        ]
                    }
                ],
                "location": "Internal",
                "defaultTags": [
                    "Element",
                    "Software System"
                ]
            },
            {
                "id": "5",
                "tags": "Element,Software System,Existing System",
                "name": "E-mail System",
                "description": "The internal Microsoft Exchange e-mail system.",
                "relationships": [
                    {
                        "id": "22",
                        "tags": "Relationship",
                        "sourceId": "5",
                        "destinationId": "1",
                        "description": "Sends e-mails to",
                        "defaultTags": [
                            "Relationship"
                        ]
                    }
                ],
                "location": "Internal",
                "defaultTags": [
                    "Element",
                    "Software System"
                ]
            },
            {
                "id": "7",
                "tags": "Element,Software System",
                "name": "Internet Banking System",
                "description": "Allows customers to view information about their bank accounts, and make payments.",
                "relationships": [
                    {
                        "id": "20",
                        "tags": "Relationship",
                        "sourceId": "7",
                        "destinationId": "4",
                        "description": "Gets account information from, and makes payments using",
                        "defaultTags": [
                            "Relationship"
                        ]
                    },
                    {
                        "id": "21",
                        "tags": "Relationship",
                        "sourceId": "7",
                        "destinationId": "5",
                        "description": "Sends e-mail using",
                        "defaultTags": [
                            "Relationship"
                        ]
                    }
                ],
                "location": "Internal",
                "containers": [
                    {
                        "id": "11",
                        "tags": "Element,Container",
                        "name": "API Application",
                        "description": "Provides Internet banking functionality via a JSON/HTTPS API.",
                        "relationships": [
                            {
                                "id": "45",
                                "sourceId": "11",
                                "destinationId": "18",
                                "description": "Reads from and writes to",
                                "technology": "JDBC",
                                "linkedRelationshipId": "44"
                            },
                            {
                                "id": "47",
                                "sourceId": "11",
                                "destinationId": "4",
                                "description": "Makes API calls to",
                                "technology": "XML/HTTPS",
                                "linkedRelationshipId": "46"
                            },
                            {
                                "id": "49",
                                "sourceId": "11",
                                "destinationId": "5",
                                "description": "Sends e-mail using",
                                "linkedRelationshipId": "48"
                            }
                        ],
                        "technology": "Java and Spring MVC",
                        "components": [
                            {
                                "id": "13",
                                "tags": "Element,Component",
                                "name": "Accounts Summary Controller",
                                "description": "Provides customers with a summary of their bank accounts.",
                                "relationships": [
                                    {
                                        "id": "41",
                                        "tags": "Relationship",
                                        "sourceId": "13",
                                        "destinationId": "16",
                                        "description": "Uses",
                                        "defaultTags": [
                                            "Relationship"
                                        ]
                                    }
                                ],
                                "technology": "Spring MVC Rest Controller",
                                "size": 0,
                                "defaultTags": [
                                    "Element",
                                    "Component"
                                ]
                            },
                            {
                                "id": "17",
                                "tags": "Element,Component",
                                "name": "E-mail Component",
                                "description": "Sends e-mails to users.",
                                "relationships": [
                                    {
                                        "id": "48",
                                        "tags": "Relationship",
                                        "sourceId": "17",
                                        "destinationId": "5",
                                        "description": "Sends e-mail using",
                                        "defaultTags": [
                                            "Relationship"
                                        ]
                                    }
                                ],
                                "technology": "Spring Bean",
                                "size": 0,
                                "defaultTags": [
                                    "Element",
                                    "Component"
                                ]
                            },
                            {
                                "id": "16",
                                "tags": "Element,Component",
                                "name": "Mainframe Banking System Facade",
                                "description": "A facade onto the mainframe banking system.",
                                "relationships": [
                                    {
                                        "id": "46",
                                        "tags": "Relationship",
                                        "sourceId": "16",
                                        "destinationId": "4",
                                        "description": "Makes API calls to",
                                        "technology": "XML/HTTPS",
                                        "defaultTags": [
                                            "Relationship"
                                        ]
                                    }
                                ],
                                "technology": "Spring Bean",
                                "size": 0,
                                "defaultTags": [
                                    "Element",
                                    "Component"
                                ]
                            },
                            {
                                "id": "14",
                                "tags": "Element,Component",
                                "name": "Reset Password Controller",
                                "description": "Allows users to reset their passwords with a single use URL.",
                                "relationships": [
                                    {
                                        "id": "42",
                                        "tags": "Relationship",
                                        "sourceId": "14",
                                        "destinationId": "15",
                                        "description": "Uses",
                                        "defaultTags": [
                                            "Relationship"
                                        ]
                                    },
                                    {
                                        "id": "43",
                                        "tags": "Relationship",
                                        "sourceId": "14",
                                        "destinationId": "17",
                                        "description": "Uses",
                                        "defaultTags": [
                                            "Relationship"
                                        ]
                                    }
                                ],
                                "technology": "Spring MVC Rest Controller",
                                "size": 0,
                                "defaultTags": [
                                    "Element",
                                    "Component"
                                ]
                            },
                            {
                                "id": "15",
                                "tags": "Element,Component",
                                "name": "Security Component",
                                "description": "Provides functionality related to signing in, changing passwords, etc.",
                                "relationships": [
                                    {
                                        "id": "44",
                                        "tags": "Relationship",
                                        "sourceId": "15",
                                        "destinationId": "18",
                                        "description": "Reads from and writes to",
                                        "technology": "JDBC",
                                        "defaultTags": [
                                            "Relationship"
                                        ]
                                    }
                                ],
                                "technology": "Spring Bean",
                                "size": 0,
                                "defaultTags": [
                                    "Element",
                                    "Component"
                                ]
                            },
                            {
                                "id": "12",
                                "tags": "Element,Component",
                                "name": "Sign In Controller",
                                "description": "Allows users to sign in to the Internet Banking System.",
                                "relationships": [
                                    {
                                        "id": "40",
                                        "tags": "Relationship",
                                        "sourceId": "12",
                                        "destinationId": "15",
                                        "description": "Uses",
                                        "defaultTags": [
                                            "Relationship"
                                        ]
                                    }
                                ],
                                "technology": "Spring MVC Rest Controller",
                                "size": 0,
                                "defaultTags": [
                                    "Element",
                                    "Component"
                                ]
                            }
                        ],
                        "defaultTags": [
                            "Element",
                            "Container"
                        ]
                    },
                    {
                        "id": "18",
                        "tags": "Element,Container,Database",
                        "name": "Database",
                        "description": "Stores user registration information, hashed authentication credentials, access logs, etc.",
                        "technology": "Oracle Database Schema",
                        "defaultTags": [
                            "Element",
                            "Container"
                        ]
                    },
                    {
                        "id": "9",
                        "tags": "Element,Container,Mobile App",
                        "name": "Mobile App",
                        "description": "Provides a limited subset of the Internet banking functionality to customers via their mobile device.",
                        "relationships": [
                            {
                                "id": "36",
                                "tags": "Relationship",
                                "sourceId": "9",
                                "destinationId": "12",
                                "description": "Makes API calls to",
                                "technology": "JSON/HTTPS",
                                "defaultTags": [
                                    "Relationship"
                                ]
                            },
                            {
                                "id": "37",
                                "sourceId": "9",
                                "destinationId": "11",
                                "description": "Makes API calls to",
                                "technology": "JSON/HTTPS",
                                "linkedRelationshipId": "36"
                            },
                            {
                                "id": "38",
                                "tags": "Relationship",
                                "sourceId": "9",
                                "destinationId": "13",
                                "description": "Makes API calls to",
                                "technology": "JSON/HTTPS",
                                "defaultTags": [
                                    "Relationship"
                                ]
                            },
                            {
                                "id": "39",
                                "tags": "Relationship",
                                "sourceId": "9",
                                "destinationId": "14",
                                "description": "Makes API calls to",
                                "technology": "JSON/HTTPS",
                                "defaultTags": [
                                    "Relationship"
                                ]
                            }
                        ],
                        "technology": "Xamarin",
                        "defaultTags": [
                            "Element",
                            "Container"
                        ]
                    },
                    {
                        "id": "8",
                        "tags": "Element,Container,Web Browser",
                        "name": "Single-Page Application",
                        "description": "Provides all of the Internet banking functionality to customers via their web browser.",
                        "relationships": [
                            {
                                "id": "32",
                                "tags": "Relationship",
                                "sourceId": "8",
                                "destinationId": "12",
                                "description": "Makes API calls to",
                                "technology": "JSON/HTTPS",
                                "defaultTags": [
                                    "Relationship"
                                ]
                            },
                            {
                                "id": "33",
                                "sourceId": "8",
                                "destinationId": "11",
                                "description": "Makes API calls to",
                                "technology": "JSON/HTTPS",
                                "linkedRelationshipId": "32"
                            },
                            {
                                "id": "34",
                                "tags": "Relationship",
                                "sourceId": "8",
                                "destinationId": "13",
                                "description": "Makes API calls to",
                                "technology": "JSON/HTTPS",
                                "defaultTags": [
                                    "Relationship"
                                ]
                            },
                            {
                                "id": "35",
                                "tags": "Relationship",
                                "sourceId": "8",
                                "destinationId": "14",
                                "description": "Makes API calls to",
                                "technology": "JSON/HTTPS",
                                "defaultTags": [
                                    "Relationship"
                                ]
                            }
                        ],
                        "technology": "JavaScript and Angular",
                        "defaultTags": [
                            "Element",
                            "Container"
                        ]
                    },
                    {
                        "id": "10",
                        "tags": "Element,Container",
                        "name": "Web Application",
                        "description": "Delivers the static content and the Internet banking single page application.",
                        "relationships": [
                            {
                                "id": "31",
                                "tags": "Relationship",
                                "sourceId": "10",
                                "destinationId": "8",
                                "description": "Delivers to the customer's web browser",
                                "defaultTags": [
                                    "Relationship"
                                ]
                            }
                        ],
                        "technology": "Java and Spring MVC",
                        "defaultTags": [
                            "Element",
                            "Container"
                        ]
                    }
                ],
                "defaultTags": [
                    "Element",
                    "Software System"
                ]
            },
            {
                "id": "4",
                "tags": "Element,Software System,Existing System",
                "name": "Mainframe Banking System",
                "description": "Stores all of the core banking information about customers, accounts, transactions, etc.",
                "location": "Internal",
                "defaultTags": [
                    "Element",
                    "Software System"
                ]
            }
        ],
        "deploymentNodes": [
            {
                "id": "63",
                "tags": "Element,Deployment Node,",
                "name": "Big Bank plc",
                "environment": "Development",
                "technology": "Big Bank plc data center",
                "instances": "1",
                "children": [
                    {
                        "id": "64",
                        "tags": "Element,Deployment Node,",
                        "name": "bigbank-dev001",
                        "environment": "Development",
                        "instances": "1",
                        "softwareSystemInstances": [
                            {
                                "id": "65",
                                "tags": "Software System Instance",
                                "environment": "Development",
                                "deploymentGroups": [
                                    "Default"
                                ],
                                "instanceId": 1,
                                "softwareSystemId": "4"
                            }
                        ],
                        "children": [],
                        "containerInstances": [],
                        "infrastructureNodes": []
                    }
                ],
                "softwareSystemInstances": [],
                "containerInstances": [],
                "infrastructureNodes": []
            },
            {
                "id": "50",
                "tags": "Element,Deployment Node",
                "name": "Developer Laptop",
                "environment": "Development",
                "technology": "Microsoft Windows 10 or Apple macOS",
                "instances": "1",
                "children": [
                    {
                        "id": "59",
                        "tags": "Element,Deployment Node",
                        "name": "Docker Container - Database Server",
                        "environment": "Development",
                        "technology": "Docker",
                        "instances": "1",
                        "children": [
                            {
                                "id": "60",
                                "tags": "Element,Deployment Node",
                                "name": "Database Server",
                                "environment": "Development",
                                "technology": "Oracle 12c",
                                "instances": "1",
                                "containerInstances": [
                                    {
                                        "id": "61",
                                        "tags": "Container Instance",
                                        "environment": "Development",
                                        "deploymentGroups": [
                                            "Default"
                                        ],
                                        "instanceId": 1,
                                        "containerId": "18"
                                    }
                                ],
                                "children": [],
                                "softwareSystemInstances": [],
                                "infrastructureNodes": []
                            }
                        ],
                        "softwareSystemInstances": [],
                        "containerInstances": [],
                        "infrastructureNodes": []
                    },
                    {
                        "id": "53",
                        "tags": "Element,Deployment Node",
                        "name": "Docker Container - Web Server",
                        "environment": "Development",
                        "technology": "Docker",
                        "instances": "1",
                        "children": [
                            {
                                "id": "54",
                                "tags": "Element,Deployment Node",
                                "name": "Apache Tomcat",
                                "environment": "Development",
                                "technology": "Apache Tomcat 8.x",
                                "instances": "1",
                                "containerInstances": [
                                    {
                                        "id": "57",
                                        "tags": "Container Instance",
                                        "relationships": [
                                            {
                                                "id": "62",
                                                "sourceId": "57",
                                                "destinationId": "61",
                                                "description": "Reads from and writes to",
                                                "technology": "JDBC",
                                                "linkedRelationshipId": "45"
                                            },
                                            {
                                                "id": "66",
                                                "sourceId": "57",
                                                "destinationId": "65",
                                                "description": "Makes API calls to",
                                                "technology": "XML/HTTPS",
                                                "linkedRelationshipId": "47"
                                            }
                                        ],
                                        "environment": "Development",
                                        "deploymentGroups": [
                                            "Default"
                                        ],
                                        "instanceId": 1,
                                        "containerId": "11"
                                    },
                                    {
                                        "id": "55",
                                        "tags": "Container Instance",
                                        "relationships": [
                                            {
                                                "id": "56",
                                                "sourceId": "55",
                                                "destinationId": "52",
                                                "description": "Delivers to the customer's web browser",
                                                "linkedRelationshipId": "31"
                                            }
                                        ],
                                        "environment": "Development",
                                        "deploymentGroups": [
                                            "Default"
                                        ],
                                        "instanceId": 1,
                                        "containerId": "10"
                                    }
                                ],
                                "children": [],
                                "softwareSystemInstances": [],
                                "infrastructureNodes": []
                            }
                        ],
                        "softwareSystemInstances": [],
                        "containerInstances": [],
                        "infrastructureNodes": []
                    },
                    {
                        "id": "51",
                        "tags": "Element,Deployment Node",
                        "name": "Web Browser",
                        "environment": "Development",
                        "technology": "Chrome, Firefox, Safari, or Edge",
                        "instances": "1",
                        "containerInstances": [
                            {
                                "id": "52",
                                "tags": "Container Instance",
                                "relationships": [
                                    {
                                        "id": "58",
                                        "sourceId": "52",
                                        "destinationId": "57",
                                        "description": "Makes API calls to",
                                        "technology": "JSON/HTTPS",
                                        "linkedRelationshipId": "33"
                                    }
                                ],
                                "environment": "Development",
                                "deploymentGroups": [
                                    "Default"
                                ],
                                "instanceId": 1,
                                "containerId": "8"
                            }
                        ],
                        "children": [],
                        "softwareSystemInstances": [],
                        "infrastructureNodes": []
                    }
                ],
                "softwareSystemInstances": [],
                "containerInstances": [],
                "infrastructureNodes": []
            },
            {
                "id": "72",
                "tags": "Element,Deployment Node",
                "name": "Big Bank plc",
                "environment": "Live",
                "technology": "Big Bank plc data center",
                "instances": "1",
                "children": [
                    {
                        "id": "77",
                        "tags": "Element,Deployment Node,",
                        "name": "bigbank-api***",
                        "environment": "Live",
                        "technology": "Ubuntu 16.04 LTS",
                        "instances": "8",
                        "children": [
                            {
                                "id": "78",
                                "tags": "Element,Deployment Node",
                                "name": "Apache Tomcat",
                                "environment": "Live",
                                "technology": "Apache Tomcat 8.x",
                                "instances": "1",
                                "containerInstances": [
                                    {
                                        "id": "79",
                                        "tags": "Container Instance",
                                        "relationships": [
                                            {
                                                "id": "85",
                                                "sourceId": "79",
                                                "destinationId": "84",
                                                "description": "Reads from and writes to",
                                                "technology": "JDBC",
                                                "linkedRelationshipId": "45"
                                            },
                                            {
                                                "id": "89",
                                                "sourceId": "79",
                                                "destinationId": "88",
                                                "description": "Reads from and writes to",
                                                "technology": "JDBC",
                                                "linkedRelationshipId": "45"
                                            },
                                            {
                                                "id": "92",
                                                "sourceId": "79",
                                                "destinationId": "91",
                                                "description": "Makes API calls to",
                                                "technology": "XML/HTTPS",
                                                "linkedRelationshipId": "47"
                                            }
                                        ],
                                        "environment": "Live",
                                        "deploymentGroups": [
                                            "Default"
                                        ],
                                        "instanceId": 1,
                                        "containerId": "11"
                                    }
                                ],
                                "children": [],
                                "softwareSystemInstances": [],
                                "infrastructureNodes": []
                            }
                        ],
                        "softwareSystemInstances": [],
                        "containerInstances": [],
                        "infrastructureNodes": []
                    },
                    {
                        "id": "82",
                        "tags": "Element,Deployment Node",
                        "name": "bigbank-db01",
                        "environment": "Live",
                        "technology": "Ubuntu 16.04 LTS",
                        "instances": "1",
                        "children": [
                            {
                                "id": "83",
                                "tags": "Element,Deployment Node",
                                "name": "Oracle - Primary",
                                "relationships": [
                                    {
                                        "id": "93",
                                        "tags": "Relationship",
                                        "sourceId": "83",
                                        "destinationId": "87",
                                        "description": "Replicates data to",
                                        "defaultTags": [
                                            "Relationship"
                                        ]
                                    }
                                ],
                                "environment": "Live",
                                "technology": "Oracle 12c",
                                "instances": "1",
                                "containerInstances": [
                                    {
                                        "id": "84",
                                        "tags": "Container Instance",
                                        "environment": "Live",
                                        "deploymentGroups": [
                                            "Default"
                                        ],
                                        "instanceId": 1,
                                        "containerId": "18"
                                    }
                                ],
                                "children": [],
                                "softwareSystemInstances": [],
                                "infrastructureNodes": []
                            }
                        ],
                        "softwareSystemInstances": [],
                        "containerInstances": [],
                        "infrastructureNodes": []
                    },
                    {
                        "id": "86",
                        "tags": "Element,Deployment Node,Failover",
                        "name": "bigbank-db02",
                        "environment": "Live",
                        "technology": "Ubuntu 16.04 LTS",
                        "instances": "1",
                        "children": [
                            {
                                "id": "87",
                                "tags": "Element,Deployment Node,Failover",
                                "name": "Oracle - Secondary",
                                "environment": "Live",
                                "technology": "Oracle 12c",
                                "instances": "1",
                                "containerInstances": [
                                    {
                                        "id": "88",
                                        "tags": "Container Instance",
                                        "environment": "Live",
                                        "deploymentGroups": [
                                            "Default"
                                        ],
                                        "instanceId": 1,
                                        "containerId": "18"
                                    }
                                ],
                                "children": [],
                                "softwareSystemInstances": [],
                                "infrastructureNodes": []
                            }
                        ],
                        "softwareSystemInstances": [],
                        "containerInstances": [],
                        "infrastructureNodes": []
                    },
                    {
                        "id": "90",
                        "tags": "Element,Deployment Node,",
                        "name": "bigbank-prod001",
                        "environment": "Live",
                        "instances": "1",
                        "softwareSystemInstances": [
                            {
                                "id": "91",
                                "tags": "Software System Instance",
                                "environment": "Live",
                                "deploymentGroups": [
                                    "Default"
                                ],
                                "instanceId": 1,
                                "softwareSystemId": "4"
                            }
                        ],
                        "children": [],
                        "containerInstances": [],
                        "infrastructureNodes": []
                    },
                    {
                        "id": "73",
                        "tags": "Element,Deployment Node,",
                        "name": "bigbank-web***",
                        "environment": "Live",
                        "technology": "Ubuntu 16.04 LTS",
                        "instances": "4",
                        "children": [
                            {
                                "id": "74",
                                "tags": "Element,Deployment Node",
                                "name": "Apache Tomcat",
                                "environment": "Live",
                                "technology": "Apache Tomcat 8.x",
                                "instances": "1",
                                "containerInstances": [
                                    {
                                        "id": "75",
                                        "tags": "Container Instance",
                                        "relationships": [
                                            {
                                                "id": "76",
                                                "sourceId": "75",
                                                "destinationId": "71",
                                                "description": "Delivers to the customer's web browser",
                                                "linkedRelationshipId": "31"
                                            }
                                        ],
                                        "environment": "Live",
                                        "deploymentGroups": [
                                            "Default"
                                        ],
                                        "instanceId": 1,
                                        "containerId": "10"
                                    }
                                ],
                                "children": [],
                                "softwareSystemInstances": [],
                                "infrastructureNodes": []
                            }
                        ],
                        "softwareSystemInstances": [],
                        "containerInstances": [],
                        "infrastructureNodes": []
                    }
                ],
                "softwareSystemInstances": [],
                "containerInstances": [],
                "infrastructureNodes": []
            },
            {
                "id": "69",
                "tags": "Element,Deployment Node",
                "name": "Customer's computer",
                "environment": "Live",
                "technology": "Microsoft Windows or Apple macOS",
                "instances": "1",
                "children": [
                    {
                        "id": "70",
                        "tags": "Element,Deployment Node",
                        "name": "Web Browser",
                        "environment": "Live",
                        "technology": "Chrome, Firefox, Safari, or Edge",
                        "instances": "1",
                        "containerInstances": [
                            {
                                "id": "71",
                                "tags": "Container Instance",
                                "relationships": [
                                    {
                                        "id": "80",
                                        "sourceId": "71",
                                        "destinationId": "79",
                                        "description": "Makes API calls to",
                                        "technology": "JSON/HTTPS",
                                        "linkedRelationshipId": "33"
                                    }
                                ],
                                "environment": "Live",
                                "deploymentGroups": [
                                    "Default"
                                ],
                                "instanceId": 1,
                                "containerId": "8"
                            }
                        ],
                        "children": [],
                        "softwareSystemInstances": [],
                        "infrastructureNodes": []
                    }
                ],
                "softwareSystemInstances": [],
                "containerInstances": [],
                "infrastructureNodes": []
            },
            {
                "id": "67",
                "tags": "Element,Deployment Node",
                "name": "Customer's mobile device",
                "environment": "Live",
                "technology": "Apple iOS or Android",
                "instances": "1",
                "containerInstances": [
                    {
                        "id": "68",
                        "tags": "Container Instance",
                        "relationships": [
                            {
                                "id": "81",
                                "sourceId": "68",
                                "destinationId": "79",
                                "description": "Makes API calls to",
                                "technology": "JSON/HTTPS",
                                "linkedRelationshipId": "37"
                            }
                        ],
                        "environment": "Live",
                        "deploymentGroups": [
                            "Default"
                        ],
                        "instanceId": 1,
                        "containerId": "9"
                    }
                ],
                "children": [],
                "softwareSystemInstances": [],
                "infrastructureNodes": []
            }
        ],
        "customElements": []
    },
    "documentation": {
        "sections": [],
        "decisions": [],
        "images": []
    },
    "views": {
        "systemLandscapeViews": [
            {
                "key": "SystemLandscape",
                "automaticLayout": {
                    "implementation": "Graphviz",
                    "rankDirection": "TopBottom",
                    "rankSeparation": 300,
                    "nodeSeparation": 300,
                    "edgeSeparation": 0,
                    "vertices": false
                },
                "enterpriseBoundaryVisible": true,
                "elements": [
                    {
                        "id": "1",
                        "x": 845,
                        "y": 213
                    },
                    {
                        "id": "2",
                        "x": 1945,
                        "y": 913
                    },
                    {
                        "id": "3",
                        "x": 2645,
                        "y": 913
                    },
                    {
                        "id": "4",
                        "x": 1558,
                        "y": 1613
                    },
                    {
                        "id": "5",
                        "x": 445,
                        "y": 1613
                    },
                    {
                        "id": "6",
                        "x": 1195,
                        "y": 963
                    },
                    {
                        "id": "7",
                        "x": 445,
                        "y": 963
                    }
                ],
                "relationships": [
                    {
                        "id": "27"
                    },
                    {
                        "id": "26"
                    },
                    {
                        "id": "25"
                    },
                    {
                        "id": "24"
                    },
                    {
                        "id": "23",
                        "vertices": [
                            {
                                "x": 1795,
                                "y": 809
                            }
                        ]
                    },
                    {
                        "id": "22",
                        "vertices": [
                            {
                                "x": 295,
                                "y": 1313
                            },
                            {
                                "x": 295,
                                "y": 809
                            }
                        ]
                    },
                    {
                        "id": "21"
                    },
                    {
                        "id": "20"
                    },
                    {
                        "id": "19"
                    }
                ],
                "animations": [],
                "paperSize": "A4_Landscape",
                "dimensions": {
                    "width": 3270,
                    "height": 2220
                }
            }
        ],
        "systemContextViews": [
            {
                "softwareSystemId": "7",
                "key": "SystemContext",
                "automaticLayout": {
                    "implementation": "Graphviz",
                    "rankDirection": "TopBottom",
                    "rankSeparation": 300,
                    "nodeSeparation": 300,
                    "edgeSeparation": 0,
                    "vertices": false
                },
                "animations": [
                    {
                        "order": 1,
                        "elements": [
                            "7"
                        ],
                        "relationships": []
                    },
                    {
                        "order": 2,
                        "elements": [
                            "1"
                        ],
                        "relationships": [
                            "19"
                        ]
                    },
                    {
                        "order": 3,
                        "elements": [
                            "4"
                        ],
                        "relationships": [
                            "20"
                        ]
                    },
                    {
                        "order": 4,
                        "elements": [
                            "5"
                        ],
                        "relationships": [
                            "22",
                            "21"
                        ]
                    }
                ],
                "enterpriseBoundaryVisible": true,
                "elements": [
                    {
                        "id": "1",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "4",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "5",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "7",
                        "x": 0,
                        "y": 0
                    }
                ],
                "relationships": [
                    {
                        "id": "22"
                    },
                    {
                        "id": "21"
                    },
                    {
                        "id": "20"
                    },
                    {
                        "id": "19"
                    }
                ]
            }
        ],
        "containerViews": [
            {
                "softwareSystemId": "7",
                "key": "Containers",
                "automaticLayout": {
                    "implementation": "Graphviz",
                    "rankDirection": "TopBottom",
                    "rankSeparation": 300,
                    "nodeSeparation": 300,
                    "edgeSeparation": 0,
                    "vertices": false
                },
                "animations": [
                    {
                        "order": 1,
                        "elements": [
                            "1",
                            "4",
                            "5"
                        ],
                        "relationships": [
                            "22"
                        ]
                    },
                    {
                        "order": 2,
                        "elements": [
                            "10"
                        ],
                        "relationships": [
                            "28"
                        ]
                    },
                    {
                        "order": 3,
                        "elements": [
                            "8"
                        ],
                        "relationships": [
                            "29",
                            "31"
                        ]
                    },
                    {
                        "order": 4,
                        "elements": [
                            "9"
                        ],
                        "relationships": [
                            "30"
                        ]
                    },
                    {
                        "order": 5,
                        "elements": [
                            "11"
                        ],
                        "relationships": [
                            "33",
                            "47",
                            "37",
                            "49"
                        ]
                    },
                    {
                        "order": 6,
                        "elements": [
                            "18"
                        ],
                        "relationships": [
                            "45"
                        ]
                    }
                ],
                "externalSoftwareSystemBoundariesVisible": true,
                "elements": [
                    {
                        "id": "11",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "1",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "4",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "5",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "18",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "8",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "9",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "10",
                        "x": 0,
                        "y": 0
                    }
                ],
                "relationships": [
                    {
                        "id": "29"
                    },
                    {
                        "id": "28"
                    },
                    {
                        "id": "37"
                    },
                    {
                        "id": "33"
                    },
                    {
                        "id": "22"
                    },
                    {
                        "id": "31"
                    },
                    {
                        "id": "45"
                    },
                    {
                        "id": "30"
                    },
                    {
                        "id": "47"
                    },
                    {
                        "id": "49"
                    }
                ]
            }
        ],
        "componentViews": [
            {
                "key": "Components",
                "automaticLayout": {
                    "implementation": "Graphviz",
                    "rankDirection": "TopBottom",
                    "rankSeparation": 300,
                    "nodeSeparation": 300,
                    "edgeSeparation": 0,
                    "vertices": false
                },
                "animations": [
                    {
                        "order": 1,
                        "elements": [
                            "4",
                            "5",
                            "18",
                            "8",
                            "9"
                        ],
                        "relationships": []
                    },
                    {
                        "order": 2,
                        "elements": [
                            "12",
                            "15"
                        ],
                        "relationships": [
                            "44",
                            "36",
                            "40",
                            "32"
                        ]
                    },
                    {
                        "order": 3,
                        "elements": [
                            "13",
                            "16"
                        ],
                        "relationships": [
                            "34",
                            "46",
                            "38",
                            "41"
                        ]
                    },
                    {
                        "order": 4,
                        "elements": [
                            "14",
                            "17"
                        ],
                        "relationships": [
                            "35",
                            "48",
                            "39",
                            "42",
                            "43"
                        ]
                    }
                ],
                "containerId": "11",
                "externalContainerBoundariesVisible": true,
                "elements": [
                    {
                        "id": "12",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "13",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "14",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "15",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "4",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "5",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "16",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "17",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "18",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "8",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "9",
                        "x": 0,
                        "y": 0
                    }
                ],
                "relationships": [
                    {
                        "id": "40"
                    },
                    {
                        "id": "41"
                    },
                    {
                        "id": "42"
                    },
                    {
                        "id": "43"
                    },
                    {
                        "id": "32"
                    },
                    {
                        "id": "36"
                    },
                    {
                        "id": "35"
                    },
                    {
                        "id": "34"
                    },
                    {
                        "id": "44"
                    },
                    {
                        "id": "46"
                    },
                    {
                        "id": "48"
                    },
                    {
                        "id": "38"
                    },
                    {
                        "id": "39"
                    }
                ]
            }
        ],
        "dynamicViews": [
            {
                "description": "Summarises how the sign in feature works in the single-page application.",
                "key": "SignIn",
                "automaticLayout": {
                    "implementation": "Graphviz",
                    "rankDirection": "TopBottom",
                    "rankSeparation": 300,
                    "nodeSeparation": 300,
                    "edgeSeparation": 0,
                    "vertices": false
                },
                "elementId": "11",
                "externalBoundariesVisible": true,
                "relationships": [
                    {
                        "id": "32",
                        "description": "Submits credentials to",
                        "order": "1",
                        "response": false
                    },
                    {
                        "id": "40",
                        "description": "Validates credentials using",
                        "order": "2",
                        "response": false
                    },
                    {
                        "id": "44",
                        "description": "select * from users where username = ?",
                        "order": "3",
                        "response": false
                    },
                    {
                        "id": "44",
                        "description": "Returns user data to",
                        "order": "4",
                        "response": true
                    },
                    {
                        "id": "40",
                        "description": "Returns true if the hashed password matches",
                        "order": "5",
                        "response": true
                    },
                    {
                        "id": "32",
                        "description": "Sends back an authentication token to",
                        "order": "6",
                        "response": true
                    }
                ],
                "elements": [
                    {
                        "id": "12",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "15",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "18",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "8",
                        "x": 0,
                        "y": 0
                    }
                ]
            }
        ],
        "deploymentViews": [
            {
                "softwareSystemId": "7",
                "key": "DevelopmentDeployment",
                "automaticLayout": {
                    "implementation": "Graphviz",
                    "rankDirection": "TopBottom",
                    "rankSeparation": 300,
                    "nodeSeparation": 300,
                    "edgeSeparation": 0,
                    "vertices": false
                },
                "environment": "Development",
                "animations": [
                    {
                        "order": 1,
                        "elements": [
                            "50",
                            "51",
                            "52"
                        ]
                    },
                    {
                        "order": 2,
                        "elements": [
                            "55",
                            "57",
                            "53",
                            "54"
                        ],
                        "relationships": [
                            "56",
                            "58"
                        ]
                    },
                    {
                        "order": 3,
                        "elements": [
                            "59",
                            "60",
                            "61"
                        ],
                        "relationships": [
                            "62"
                        ]
                    }
                ],
                "elements": [
                    {
                        "id": "55",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "57",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "59",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "60",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "61",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "50",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "51",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "52",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "63",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "53",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "64",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "54",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "65",
                        "x": 0,
                        "y": 0
                    }
                ],
                "relationships": [
                    {
                        "id": "62"
                    },
                    {
                        "id": "56"
                    },
                    {
                        "id": "66"
                    },
                    {
                        "id": "58"
                    }
                ]
            },
            {
                "softwareSystemId": "7",
                "key": "LiveDeployment",
                "automaticLayout": {
                    "implementation": "Graphviz",
                    "rankDirection": "TopBottom",
                    "rankSeparation": 300,
                    "nodeSeparation": 300,
                    "edgeSeparation": 0,
                    "vertices": false
                },
                "environment": "Live",
                "animations": [
                    {
                        "order": 1,
                        "elements": [
                            "69",
                            "70",
                            "71"
                        ]
                    },
                    {
                        "order": 2,
                        "elements": [
                            "67",
                            "68"
                        ]
                    },
                    {
                        "order": 3,
                        "elements": [
                            "77",
                            "78",
                            "79",
                            "72",
                            "73",
                            "74",
                            "75"
                        ],
                        "relationships": [
                            "80",
                            "81",
                            "76"
                        ]
                    },
                    {
                        "order": 4,
                        "elements": [
                            "82",
                            "83",
                            "84"
                        ],
                        "relationships": [
                            "85"
                        ]
                    },
                    {
                        "order": 5,
                        "elements": [
                            "88",
                            "86",
                            "87"
                        ],
                        "relationships": [
                            "89",
                            "93"
                        ]
                    }
                ],
                "elements": [
                    {
                        "id": "88",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "77",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "67",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "78",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "68",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "79",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "69",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "90",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "91",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "70",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "71",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "82",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "83",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "72",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "84",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "73",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "74",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "86",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "75",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "87",
                        "x": 0,
                        "y": 0
                    }
                ],
                "relationships": [
                    {
                        "id": "93"
                    },
                    {
                        "id": "80"
                    },
                    {
                        "id": "81"
                    },
                    {
                        "id": "92"
                    },
                    {
                        "id": "76"
                    },
                    {
                        "id": "85"
                    },
                    {
                        "id": "89"
                    }
                ]
            }
        ],
        "configuration": {
            "branding": {},
            "styles": {
                "elements": [
                    {
                        "tag": "Person",
                        "color": "#ffffff",
                        "fontSize": 22,
                        "shape": "Person"
                    },
                    {
                        "tag": "Customer",
                        "background": "#08427b"
                    },
                    {
                        "tag": "Bank Staff",
                        "background": "#999999"
                    },
                    {
                        "tag": "Software System",
                        "background": "#1168bd",
                        "color": "#ffffff"
                    },
                    {
                        "tag": "Existing System",
                        "background": "#999999",
                        "color": "#ffffff"
                    },
                    {
                        "tag": "Container",
                        "background": "#438dd5",
                        "color": "#ffffff"
                    },
                    {
                        "tag": "Web Browser",
                        "shape": "WebBrowser"
                    },
                    {
                        "tag": "Mobile App",
                        "shape": "MobileDeviceLandscape"
                    },
                    {
                        "tag": "Database",
                        "shape": "Cylinder"
                    },
                    {
                        "tag": "Component",
                        "background": "#85bbf0",
                        "color": "#000000"
                    },
                    {
                        "tag": "Failover",
                        "opacity": 25
                    }
                ],
                "relationships": []
            },
            "terminology": {},
            "themes": []
        },
        "customViews": [],
        "filteredViews": []
    }
}
`
