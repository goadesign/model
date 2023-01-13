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
		json.NewEncoder(rw).Encode(wkspc)
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
// https://raw.githubusercontent.com/structurizr/json/master/examples/big-bank-plc.json
var bigBankPLC = `{
    "name": "Big Bank plc",
    "description": "This is an example workspace to illustrate the key features of Structurizr, based around a fictional online banking system.",
    "model": {
        "enterprise": {
            "name": "Big Bank plc"
        },
        "people": [
            {
                "id": "15",
                "tags": "Element,Person,Bank Staff",
                "name": "Back Office Staff",
                "description": "Administration and support staff within the bank.",
                "relationships": [
                    {
                        "id": "16",
                        "tags": "Relationship,Synchronous",
                        "sourceId": "15",
                        "destinationId": "4",
                        "description": "Uses",
                        "interactionStyle": "Synchronous"
                    }
                ],
                "location": "Internal"
            },
            {
                "id": "12",
                "tags": "Element,Person,Bank Staff",
                "name": "Customer Service Staff",
                "description": "Customer service staff within the bank.",
                "relationships": [
                    {
                        "id": "13",
                        "tags": "Relationship,Synchronous",
                        "sourceId": "12",
                        "destinationId": "4",
                        "description": "Uses",
                        "interactionStyle": "Synchronous"
                    }
                ],
                "location": "Internal"
            },
            {
                "id": "1",
                "tags": "Element,Person",
                "name": "Personal Banking Customer",
                "description": "A customer of the bank, with personal bank accounts.",
                "relationships": [
                    {
                        "id": "23",
                        "tags": "Relationship,Synchronous",
                        "sourceId": "1",
                        "destinationId": "17",
                        "description": "Views account balances, and makes payments using",
                        "interactionStyle": "Synchronous"
                    },
                    {
                        "id": "11",
                        "tags": "Relationship,Synchronous",
                        "sourceId": "1",
                        "destinationId": "9",
                        "description": "Withdraws cash using",
                        "interactionStyle": "Synchronous"
                    },
                    {
                        "id": "14",
                        "tags": "Relationship,Synchronous",
                        "sourceId": "1",
                        "destinationId": "12",
                        "description": "Asks questions to",
                        "technology": "Telephone",
                        "interactionStyle": "Synchronous"
                    },
                    {
                        "id": "3",
                        "tags": "Relationship,Synchronous",
                        "sourceId": "1",
                        "destinationId": "2",
                        "description": "Views account balances, and makes payments using",
                        "interactionStyle": "Synchronous"
                    },
                    {
                        "id": "24",
                        "tags": "Relationship,Synchronous",
                        "sourceId": "1",
                        "destinationId": "18",
                        "description": "Views account balances, and makes payments using",
                        "interactionStyle": "Synchronous"
                    },
                    {
                        "id": "22",
                        "tags": "Relationship,Synchronous",
                        "sourceId": "1",
                        "destinationId": "19",
                        "description": "Visits bigbank.com/ib using",
                        "technology": "HTTPS",
                        "interactionStyle": "Synchronous"
                    }
                ],
                "location": "External"
            }
        ],
        "softwareSystems": [
            {
                "id": "9",
                "tags": "Element,Software System,Existing System",
                "name": "ATM",
                "description": "Allows customers to withdraw cash.",
                "relationships": [
                    {
                        "id": "10",
                        "tags": "Relationship,Synchronous",
                        "sourceId": "9",
                        "destinationId": "4",
                        "description": "Uses",
                        "interactionStyle": "Synchronous"
                    }
                ],
                "location": "Internal"
            },
            {
                "id": "6",
                "tags": "Element,Software System,Existing System",
                "name": "E-mail System",
                "description": "The internal Microsoft Exchange e-mail system.",
                "relationships": [
                    {
                        "id": "8",
                        "tags": "Relationship,Synchronous",
                        "sourceId": "6",
                        "destinationId": "1",
                        "description": "Sends e-mails to",
                        "interactionStyle": "Synchronous"
                    }
                ],
                "location": "Internal"
            },
            {
                "id": "2",
                "tags": "Element,Software System",
                "name": "Internet Banking System",
                "description": "Allows customers to view information about their bank accounts, and make payments.",
                "relationships": [
                    {
                        "id": "7",
                        "tags": "Relationship,Synchronous",
                        "sourceId": "2",
                        "destinationId": "6",
                        "description": "Sends e-mail using",
                        "interactionStyle": "Synchronous"
                    },
                    {
                        "id": "5",
                        "tags": "Relationship,Synchronous",
                        "sourceId": "2",
                        "destinationId": "4",
                        "description": "Gets account information from, and makes payments using",
                        "interactionStyle": "Synchronous"
                    }
                ],
                "location": "Internal",
                "containers": [
                    {
                        "id": "20",
                        "tags": "Element,Container",
                        "name": "API Application",
                        "description": "Provides Internet banking functionality via a JSON/HTTPS API.",
                        "relationships": [
                            {
                                "id": "27",
                                "tags": "Relationship,Synchronous",
                                "sourceId": "20",
                                "destinationId": "4",
                                "description": "Makes API calls to",
                                "technology": "XML/HTTPS",
                                "interactionStyle": "Synchronous"
                            },
                            {
                                "id": "26",
                                "tags": "Relationship,Synchronous",
                                "sourceId": "20",
                                "destinationId": "21",
                                "description": "Reads from and writes to",
                                "technology": "JDBC",
                                "interactionStyle": "Synchronous"
                            },
                            {
                                "id": "28",
                                "tags": "Relationship,Synchronous",
                                "sourceId": "20",
                                "destinationId": "6",
                                "description": "Sends e-mail using",
                                "technology": "SMTP",
                                "interactionStyle": "Synchronous"
                            }
                        ],
                        "technology": "Java and Spring MVC",
                        "components": [
                            {
                                "id": "30",
                                "tags": "Element,Component",
                                "name": "Accounts Summary Controller",
                                "description": "Provides customers with a summary of their bank accounts.",
                                "relationships": [
                                    {
                                        "id": "42",
                                        "tags": "Relationship,Synchronous",
                                        "sourceId": "30",
                                        "destinationId": "33",
                                        "description": "Uses",
                                        "interactionStyle": "Synchronous"
                                    }
                                ],
                                "technology": "Spring MVC Rest Controller",
                                "size": 0
                            },
                            {
                                "id": "34",
                                "tags": "Element,Component",
                                "name": "E-mail Component",
                                "description": "Sends e-mails to users.",
                                "relationships": [
                                    {
                                        "id": "47",
                                        "tags": "Relationship,Synchronous",
                                        "sourceId": "34",
                                        "destinationId": "6",
                                        "description": "Sends e-mail using",
                                        "interactionStyle": "Synchronous"
                                    }
                                ],
                                "technology": "Spring Bean",
                                "size": 0
                            },
                            {
                                "id": "33",
                                "tags": "Element,Component",
                                "name": "Mainframe Banking System Facade",
                                "description": "A facade onto the mainframe banking system.",
                                "relationships": [
                                    {
                                        "id": "46",
                                        "tags": "Relationship,Synchronous",
                                        "sourceId": "33",
                                        "destinationId": "4",
                                        "description": "Uses",
                                        "technology": "XML/HTTPS",
                                        "interactionStyle": "Synchronous"
                                    }
                                ],
                                "technology": "Spring Bean",
                                "size": 0
                            },
                            {
                                "id": "31",
                                "tags": "Element,Component",
                                "name": "Reset Password Controller",
                                "description": "Allows users to reset their passwords with a single use URL.",
                                "relationships": [
                                    {
                                        "id": "44",
                                        "tags": "Relationship,Synchronous",
                                        "sourceId": "31",
                                        "destinationId": "34",
                                        "description": "Uses",
                                        "interactionStyle": "Synchronous"
                                    },
                                    {
                                        "id": "43",
                                        "tags": "Relationship,Synchronous",
                                        "sourceId": "31",
                                        "destinationId": "32",
                                        "description": "Uses",
                                        "interactionStyle": "Synchronous"
                                    }
                                ],
                                "technology": "Spring MVC Rest Controller",
                                "size": 0
                            },
                            {
                                "id": "32",
                                "tags": "Element,Component",
                                "name": "Security Component",
                                "description": "Provides functionality related to signing in, changing passwords, etc.",
                                "relationships": [
                                    {
                                        "id": "45",
                                        "tags": "Relationship,Synchronous",
                                        "sourceId": "32",
                                        "destinationId": "21",
                                        "description": "Reads from and writes to",
                                        "technology": "JDBC",
                                        "interactionStyle": "Synchronous"
                                    }
                                ],
                                "technology": "Spring Bean",
                                "size": 0
                            },
                            {
                                "id": "29",
                                "tags": "Element,Component",
                                "name": "Sign In Controller",
                                "description": "Allows users to sign in to the Internet Banking System.",
                                "relationships": [
                                    {
                                        "id": "41",
                                        "tags": "Relationship,Synchronous",
                                        "sourceId": "29",
                                        "destinationId": "32",
                                        "description": "Uses",
                                        "interactionStyle": "Synchronous"
                                    }
                                ],
                                "technology": "Spring MVC Rest Controller",
                                "size": 0
                            }
                        ]
                    },
                    {
                        "id": "21",
                        "tags": "Element,Container,Database",
                        "name": "Database",
                        "description": "Stores user registration information, hashed authentication credentials, access logs, etc.",
                        "technology": "Oracle Database Schema"
                    },
                    {
                        "id": "18",
                        "tags": "Element,Container,Mobile App",
                        "name": "Mobile App",
                        "description": "Provides a limited subset of the Internet banking functionality to customers via their mobile device.",
                        "relationships": [
                            {
                                "id": "39",
                                "tags": "Relationship,Synchronous",
                                "sourceId": "18",
                                "destinationId": "31",
                                "description": "Makes API calls to",
                                "technology": "JSON/HTTPS",
                                "interactionStyle": "Synchronous"
                            },
                            {
                                "id": "49",
                                "tags": "Relationship,Synchronous",
                                "sourceId": "18",
                                "destinationId": "20",
                                "description": "Makes API calls to",
                                "technology": "JSON/HTTPS",
                                "interactionStyle": "Synchronous"
                            },
                            {
                                "id": "38",
                                "tags": "Relationship,Synchronous",
                                "sourceId": "18",
                                "destinationId": "29",
                                "description": "Makes API calls to",
                                "technology": "JSON/HTTPS",
                                "interactionStyle": "Synchronous"
                            },
                            {
                                "id": "40",
                                "tags": "Relationship,Synchronous",
                                "sourceId": "18",
                                "destinationId": "30",
                                "description": "Makes API calls to",
                                "technology": "JSON/HTTPS",
                                "interactionStyle": "Synchronous"
                            }
                        ],
                        "technology": "Xamarin"
                    },
                    {
                        "id": "17",
                        "tags": "Element,Container,Web Browser",
                        "name": "Single-Page Application",
                        "description": "Provides all of the Internet banking functionality to customers via their web browser.",
                        "relationships": [
                            {
                                "id": "37",
                                "tags": "Relationship,Synchronous",
                                "sourceId": "17",
                                "destinationId": "30",
                                "description": "Makes API calls to",
                                "technology": "JSON/HTTPS",
                                "interactionStyle": "Synchronous"
                            },
                            {
                                "id": "35",
                                "tags": "Relationship,Synchronous",
                                "sourceId": "17",
                                "destinationId": "29",
                                "description": "Makes API calls to",
                                "technology": "JSON/HTTPS",
                                "interactionStyle": "Synchronous"
                            },
                            {
                                "id": "48",
                                "tags": "Relationship,Synchronous",
                                "sourceId": "17",
                                "destinationId": "20",
                                "description": "Makes API calls to",
                                "technology": "JSON/HTTPS",
                                "interactionStyle": "Synchronous"
                            },
                            {
                                "id": "36",
                                "tags": "Relationship,Synchronous",
                                "sourceId": "17",
                                "destinationId": "31",
                                "description": "Makes API calls to",
                                "technology": "JSON/HTTPS",
                                "interactionStyle": "Synchronous"
                            }
                        ],
                        "technology": "JavaScript and Angular"
                    },
                    {
                        "id": "19",
                        "tags": "Element,Container",
                        "name": "Web Application",
                        "description": "Delivers the static content and the Internet banking single page application.",
                        "relationships": [
                            {
                                "id": "25",
                                "tags": "Relationship,Synchronous",
                                "sourceId": "19",
                                "destinationId": "17",
                                "description": "Delivers to the customer's web browser",
                                "interactionStyle": "Synchronous"
                            }
                        ],
                        "technology": "Java and Spring MVC"
                    }
                ]
            },
            {
                "id": "4",
                "tags": "Element,Software System,Existing System",
                "name": "Mainframe Banking System",
                "description": "Stores all of the core banking information about customers, accounts, transactions, etc.",
                "location": "Internal"
            }
        ],
        "deploymentNodes": [
            {
                "id": "50",
                "tags": "Element,Deployment Node",
                "name": "Developer Laptop",
                "description": "A developer laptop.",
                "environment": "Development",
                "technology": "Microsoft Windows 10 or Apple macOS",
                "instances": "1",
                "children": [
                    {
                        "id": "55",
                        "tags": "Element,Deployment Node",
                        "name": "Docker Container - Database Server",
                        "description": "A Docker container.",
                        "environment": "Development",
                        "technology": "Docker",
                        "instances": "1",
                        "children": [
                            {
                                "id": "56",
                                "tags": "Element,Deployment Node",
                                "name": "Database Server",
                                "description": "A development database.",
                                "environment": "Development",
                                "technology": "Oracle 12c",
                                "instances": "1",
                                "containerInstances": [
                                    {
                                        "id": "57",
                                        "tags": "Container Instance",
                                        "environment": "Development",
                                        "containerId": "21",
                                        "instanceId": 1,
                                        "properties": {}
                                    }
                                ],
                                "children": [],
                                "infrastructureNodes": []
                            }
                        ],
                        "containerInstances": [],
                        "infrastructureNodes": []
                    },
                    {
                        "id": "51",
                        "tags": "Element,Deployment Node",
                        "name": "Docker Container - Web Server",
                        "description": "A Docker container.",
                        "environment": "Development",
                        "technology": "Docker",
                        "instances": "1",
                        "children": [
                            {
                                "id": "52",
                                "tags": "Element,Deployment Node",
                                "properties": {
                                    "Java Version": "8",
                                    "Xms": "1024M",
                                    "Xmx": "512M"
                                },
                                "name": "Apache Tomcat",
                                "description": "An open source Java EE web server.",
                                "environment": "Development",
                                "technology": "Apache Tomcat 8.x",
                                "instances": "1",
                                "containerInstances": [
                                    {
                                        "id": "54",
                                        "tags": "Container Instance",
                                        "relationships": [
                                            {
                                                "id": "58",
                                                "sourceId": "54",
                                                "destinationId": "57",
                                                "description": "Reads from and writes to",
                                                "technology": "JDBC",
                                                "interactionStyle": "Synchronous",
                                                "linkedRelationshipId": "26"
                                            }
                                        ],
                                        "environment": "Development",
                                        "containerId": "20",
                                        "instanceId": 1,
                                        "properties": {}
                                    },
                                    {
                                        "id": "53",
                                        "tags": "Container Instance",
                                        "relationships": [
                                            {
                                                "id": "62",
                                                "sourceId": "53",
                                                "destinationId": "60",
                                                "description": "Delivers to the customer's web browser",
                                                "interactionStyle": "Synchronous",
                                                "linkedRelationshipId": "25"
                                            }
                                        ],
                                        "environment": "Development",
                                        "containerId": "19",
                                        "instanceId": 1,
                                        "properties": {}
                                    }
                                ],
                                "children": [],
                                "infrastructureNodes": []
                            }
                        ],
                        "containerInstances": [],
                        "infrastructureNodes": []
                    },
                    {
                        "id": "59",
                        "tags": "Element,Deployment Node",
                        "name": "Web Browser",
                        "environment": "Development",
                        "technology": "Chrome, Firefox, Safari, or Edge",
                        "instances": "1",
                        "containerInstances": [
                            {
                                "id": "60",
                                "tags": "Container Instance",
                                "relationships": [
                                    {
                                        "id": "61",
                                        "sourceId": "60",
                                        "destinationId": "54",
                                        "description": "Makes API calls to",
                                        "technology": "JSON/HTTPS",
                                        "interactionStyle": "Synchronous",
                                        "linkedRelationshipId": "48"
                                    }
                                ],
                                "environment": "Development",
                                "containerId": "17",
                                "instanceId": 1,
                                "properties": {}
                            }
                        ],
                        "children": [],
                        "infrastructureNodes": []
                    }
                ],
                "containerInstances": [],
                "infrastructureNodes": []
            },
            {
                "id": "68",
                "tags": "Element,Deployment Node",
                "name": "Big Bank plc",
                "environment": "Live",
                "technology": "Big Bank plc data center",
                "instances": "1",
                "children": [
                    {
                        "id": "73",
                        "tags": "Element,Deployment Node",
                        "properties": {
                            "Location": "London and Reading"
                        },
                        "name": "bigbank-api***",
                        "description": "A web server residing in the web server farm, accessed via F5 BIG-IP LTMs.",
                        "environment": "Live",
                        "technology": "Ubuntu 16.04 LTS",
                        "instances": "8",
                        "children": [
                            {
                                "id": "74",
                                "tags": "Element,Deployment Node",
                                "properties": {
                                    "Java Version": "8",
                                    "Xms": "1024M",
                                    "Xmx": "512M"
                                },
                                "name": "Apache Tomcat",
                                "description": "An open source Java EE web server.",
                                "environment": "Live",
                                "technology": "Apache Tomcat 8.x",
                                "instances": "1",
                                "containerInstances": [
                                    {
                                        "id": "75",
                                        "tags": "Container Instance",
                                        "relationships": [
                                            {
                                                "id": "81",
                                                "sourceId": "75",
                                                "destinationId": "80",
                                                "description": "Reads from and writes to",
                                                "technology": "JDBC",
                                                "interactionStyle": "Synchronous",
                                                "linkedRelationshipId": "26"
                                            },
                                            {
                                                "id": "85",
                                                "tags": "Failover",
                                                "sourceId": "75",
                                                "destinationId": "84",
                                                "description": "Reads from and writes to",
                                                "technology": "JDBC",
                                                "interactionStyle": "Synchronous",
                                                "linkedRelationshipId": "26"
                                            }
                                        ],
                                        "environment": "Live",
                                        "containerId": "20",
                                        "instanceId": 2,
                                        "properties": {}
                                    }
                                ],
                                "children": [],
                                "infrastructureNodes": []
                            }
                        ],
                        "containerInstances": [],
                        "infrastructureNodes": []
                    },
                    {
                        "id": "78",
                        "tags": "Element,Deployment Node",
                        "properties": {
                            "Location": "London"
                        },
                        "name": "bigbank-db01",
                        "description": "The primary database server.",
                        "environment": "Live",
                        "technology": "Ubuntu 16.04 LTS",
                        "instances": "1",
                        "children": [
                            {
                                "id": "79",
                                "tags": "Element,Deployment Node",
                                "name": "Oracle - Primary",
                                "description": "The primary, live database server.",
                                "relationships": [
                                    {
                                        "id": "86",
                                        "tags": "Relationship,Synchronous",
                                        "sourceId": "79",
                                        "destinationId": "83",
                                        "description": "Replicates data to",
                                        "interactionStyle": "Synchronous"
                                    }
                                ],
                                "environment": "Live",
                                "technology": "Oracle 12c",
                                "instances": "1",
                                "containerInstances": [
                                    {
                                        "id": "80",
                                        "tags": "Container Instance",
                                        "environment": "Live",
                                        "containerId": "21",
                                        "instanceId": 2,
                                        "properties": {}
                                    }
                                ],
                                "children": [],
                                "infrastructureNodes": []
                            }
                        ],
                        "containerInstances": [],
                        "infrastructureNodes": []
                    },
                    {
                        "id": "82",
                        "tags": "Element,Deployment Node,Failover",
                        "properties": {
                            "Location": "Reading"
                        },
                        "name": "bigbank-db02",
                        "description": "The secondary database server.",
                        "environment": "Live",
                        "technology": "Ubuntu 16.04 LTS",
                        "instances": "1",
                        "children": [
                            {
                                "id": "83",
                                "tags": "Element,Deployment Node,Failover",
                                "name": "Oracle - Secondary",
                                "description": "A secondary, standby database server, used for failover purposes only.",
                                "environment": "Live",
                                "technology": "Oracle 12c",
                                "instances": "1",
                                "containerInstances": [
                                    {
                                        "id": "84",
                                        "tags": "Container Instance,Failover",
                                        "environment": "Live",
                                        "containerId": "21",
                                        "instanceId": 3,
                                        "properties": {}
                                    }
                                ],
                                "children": [],
                                "infrastructureNodes": []
                            }
                        ],
                        "containerInstances": [],
                        "infrastructureNodes": []
                    },
                    {
                        "id": "69",
                        "tags": "Element,Deployment Node",
                        "properties": {
                            "Location": "London and Reading"
                        },
                        "name": "bigbank-web***",
                        "description": "A web server residing in the web server farm, accessed via F5 BIG-IP LTMs.",
                        "environment": "Live",
                        "technology": "Ubuntu 16.04 LTS",
                        "instances": "4",
                        "children": [
                            {
                                "id": "70",
                                "tags": "Element,Deployment Node",
                                "properties": {
                                    "Java Version": "8",
                                    "Xms": "1024M",
                                    "Xmx": "512M"
                                },
                                "name": "Apache Tomcat",
                                "description": "An open source Java EE web server.",
                                "environment": "Live",
                                "technology": "Apache Tomcat 8.x",
                                "instances": "1",
                                "containerInstances": [
                                    {
                                        "id": "71",
                                        "tags": "Container Instance",
                                        "relationships": [
                                            {
                                                "id": "72",
                                                "sourceId": "71",
                                                "destinationId": "67",
                                                "description": "Delivers to the customer's web browser",
                                                "interactionStyle": "Synchronous",
                                                "linkedRelationshipId": "25"
                                            }
                                        ],
                                        "environment": "Live",
                                        "containerId": "19",
                                        "instanceId": 2,
                                        "properties": {}
                                    }
                                ],
                                "children": [],
                                "infrastructureNodes": []
                            }
                        ],
                        "containerInstances": [],
                        "infrastructureNodes": []
                    }
                ],
                "containerInstances": [],
                "infrastructureNodes": []
            },
            {
                "id": "65",
                "tags": "Element,Deployment Node",
                "name": "Customer's computer",
                "environment": "Live",
                "technology": "Microsoft Windows or Apple macOS",
                "instances": "1",
                "children": [
                    {
                        "id": "66",
                        "tags": "Element,Deployment Node",
                        "name": "Web Browser",
                        "environment": "Live",
                        "technology": "Chrome, Firefox, Safari, or Edge",
                        "instances": "1",
                        "containerInstances": [
                            {
                                "id": "67",
                                "tags": "Container Instance",
                                "relationships": [
                                    {
                                        "id": "77",
                                        "sourceId": "67",
                                        "destinationId": "75",
                                        "description": "Makes API calls to",
                                        "technology": "JSON/HTTPS",
                                        "interactionStyle": "Synchronous",
                                        "linkedRelationshipId": "48"
                                    }
                                ],
                                "environment": "Live",
                                "containerId": "17",
                                "instanceId": 2,
                                "properties": {}
                            }
                        ],
                        "children": [],
                        "infrastructureNodes": []
                    }
                ],
                "containerInstances": [],
                "infrastructureNodes": []
            },
            {
                "id": "63",
                "tags": "Element,Deployment Node",
                "name": "Customer's mobile device",
                "environment": "Live",
                "technology": "Apple iOS or Android",
                "instances": "1",
                "containerInstances": [
                    {
                        "id": "64",
                        "tags": "Container Instance",
                        "relationships": [
                            {
                                "id": "76",
                                "sourceId": "64",
                                "destinationId": "75",
                                "description": "Makes API calls to",
                                "technology": "JSON/HTTPS",
                                "interactionStyle": "Synchronous",
                                "linkedRelationshipId": "49"
                            }
                        ],
                        "environment": "Live",
                        "containerId": "18",
                        "instanceId": 1,
                        "properties": {}
                    }
                ],
                "children": [],
                "infrastructureNodes": []
            }
        ]
    },
    "documentation": {
        "sections": [
            {
                "elementId": "2",
                "title": "Context",
                "order": 1,
                "format": "Markdown",
                "content": "Here is some context about the Internet Banking System...\n![](embed:SystemLandscape)\n![](embed:SystemContext)\n### Internet Banking System\n...\n### Mainframe Banking System\n...\n"
            },
            {
                "elementId": "19",
                "title": "Components",
                "order": 3,
                "format": "Markdown",
                "content": "Here is some information about the API Application...\n![](embed:Components)\n### Sign in process\nHere is some information about the Sign In Controller, including how the sign in process works...\n![](embed:SignIn)"
            },
            {
                "elementId": "2",
                "title": "Development Environment",
                "order": 4,
                "format": "AsciiDoc",
                "content": "Here is some information about how to set up a development environment for the Internet Banking System...\nimage::embed:DevelopmentDeployment[]"
            },
            {
                "elementId": "2",
                "title": "Containers",
                "order": 2,
                "format": "Markdown",
                "content": "Here is some information about the containers within the Internet Banking System...\n![](embed:Containers)\n### Web Application\n...\n### Database\n...\n"
            },
            {
                "elementId": "2",
                "title": "Deployment",
                "order": 5,
                "format": "AsciiDoc",
                "content": "Here is some information about the live deployment environment for the Internet Banking System...\nimage::embed:LiveDeployment[]"
            }
        ],
        "template": {
            "name": "Software Guidebook",
            "author": "Simon Brown",
            "url": "https://leanpub.com/visualising-software-architecture"
        },
        "decisions": [],
        "images": []
    },
    "views": {
        "systemLandscapeViews": [
            {
                "description": "The system landscape diagram for Big Bank plc.",
                "key": "SystemLandscape",
                "paperSize": "A5_Landscape",
                "animations": [
                    {
                        "order": 1,
                        "elements": [
                            "1",
                            "2",
                            "4",
                            "6"
                        ],
                        "relationships": [
                            "3",
                            "5",
                            "7",
                            "8"
                        ]
                    },
                    {
                        "order": 2,
                        "elements": [
                            "9"
                        ],
                        "relationships": [
                            "11",
                            "10"
                        ]
                    },
                    {
                        "order": 3,
                        "elements": [
                            "12",
                            "15"
                        ],
                        "relationships": [
                            "13",
                            "14",
                            "16"
                        ]
                    }
                ],
                "enterpriseBoundaryVisible": true,
                "elements": [
                    {
                        "id": "1",
                        "x": 87,
                        "y": 643
                    },
                    {
                        "id": "12",
                        "x": 1947,
                        "y": 36
                    },
                    {
                        "id": "2",
                        "x": 1012,
                        "y": 813
                    },
                    {
                        "id": "4",
                        "x": 1922,
                        "y": 693
                    },
                    {
                        "id": "15",
                        "x": 1947,
                        "y": 1241
                    },
                    {
                        "id": "6",
                        "x": 1012,
                        "y": 1326
                    },
                    {
                        "id": "9",
                        "x": 1012,
                        "y": 301
                    }
                ],
                "relationships": [
                    {
                        "id": "16"
                    },
                    {
                        "id": "3"
                    },
                    {
                        "id": "14",
                        "vertices": [
                            {
                                "x": 285,
                                "y": 240
                            }
                        ]
                    },
                    {
                        "id": "5"
                    },
                    {
                        "id": "13"
                    },
                    {
                        "id": "11"
                    },
                    {
                        "id": "7"
                    },
                    {
                        "id": "8"
                    },
                    {
                        "id": "10"
                    }
                ]
            }
        ],
        "systemContextViews": [
            {
                "softwareSystemId": "2",
                "description": "The system context diagram for the Internet Banking System.",
                "key": "SystemContext",
                "paperSize": "A5_Landscape",
                "animations": [
                    {
                        "order": 1,
                        "elements": [
                            "2"
                        ],
                        "relationships": []
                    },
                    {
                        "order": 2,
                        "elements": [
                            "1"
                        ],
                        "relationships": [
                            "3"
                        ]
                    },
                    {
                        "order": 3,
                        "elements": [
                            "4"
                        ],
                        "relationships": [
                            "5"
                        ]
                    },
                    {
                        "order": 4,
                        "elements": [
                            "6"
                        ],
                        "relationships": [
                            "7",
                            "8"
                        ]
                    }
                ],
                "enterpriseBoundaryVisible": false,
                "elements": [
                    {
                        "id": "1",
                        "x": 632,
                        "y": 69
                    },
                    {
                        "id": "2",
                        "x": 607,
                        "y": 714
                    },
                    {
                        "id": "4",
                        "x": 607,
                        "y": 1259
                    },
                    {
                        "id": "6",
                        "x": 1422,
                        "y": 714
                    }
                ],
                "relationships": [
                    {
                        "id": "3"
                    },
                    {
                        "id": "5"
                    },
                    {
                        "id": "7"
                    },
                    {
                        "id": "8"
                    }
                ]
            }
        ],
        "containerViews": [
            {
                "softwareSystemId": "2",
                "description": "The container diagram for the Internet Banking System.",
                "key": "Containers",
                "paperSize": "A5_Landscape",
                "animations": [
                    {
                        "order": 1,
                        "elements": [
                            "1",
                            "4",
                            "6"
                        ],
                        "relationships": [
                            "8"
                        ]
                    },
                    {
                        "order": 2,
                        "elements": [
                            "19"
                        ],
                        "relationships": [
                            "22"
                        ]
                    },
                    {
                        "order": 3,
                        "elements": [
                            "17"
                        ],
                        "relationships": [
                            "23",
                            "25"
                        ]
                    },
                    {
                        "order": 4,
                        "elements": [
                            "18"
                        ],
                        "relationships": [
                            "24"
                        ]
                    },
                    {
                        "order": 5,
                        "elements": [
                            "20"
                        ],
                        "relationships": [
                            "48",
                            "27",
                            "49",
                            "28"
                        ]
                    },
                    {
                        "order": 6,
                        "elements": [
                            "21"
                        ],
                        "relationships": [
                            "26"
                        ]
                    }
                ],
                "externalSoftwareSystemBoundariesVisible": false,
                "elements": [
                    {
                        "id": "1",
                        "x": 1056,
                        "y": 24
                    },
                    {
                        "id": "4",
                        "x": 2012,
                        "y": 1214
                    },
                    {
                        "id": "17",
                        "x": 780,
                        "y": 664
                    },
                    {
                        "id": "6",
                        "x": 2012,
                        "y": 664
                    },
                    {
                        "id": "18",
                        "x": 1283,
                        "y": 664
                    },
                    {
                        "id": "19",
                        "x": 37,
                        "y": 664
                    },
                    {
                        "id": "20",
                        "x": 1031,
                        "y": 1214
                    },
                    {
                        "id": "21",
                        "x": 37,
                        "y": 1214
                    }
                ],
                "relationships": [
                    {
                        "id": "28"
                    },
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
                        "id": "23"
                    },
                    {
                        "id": "22"
                    },
                    {
                        "id": "8"
                    },
                    {
                        "id": "48"
                    },
                    {
                        "id": "49"
                    }
                ]
            }
        ],
        "componentViews": [
            {
                "description": "The component diagram for the API Application.",
                "key": "Components",
                "paperSize": "A5_Landscape",
                "animations": [
                    {
                        "order": 1,
                        "elements": [
                            "4",
                            "17",
                            "6",
                            "18",
                            "21"
                        ],
                        "relationships": []
                    },
                    {
                        "order": 2,
                        "elements": [
                            "29",
                            "32"
                        ],
                        "relationships": [
                            "45",
                            "35",
                            "38",
                            "41"
                        ]
                    },
                    {
                        "order": 3,
                        "elements": [
                            "33",
                            "30"
                        ],
                        "relationships": [
                            "46",
                            "37",
                            "40",
                            "42"
                        ]
                    },
                    {
                        "order": 4,
                        "elements": [
                            "34",
                            "31"
                        ],
                        "relationships": [
                            "44",
                            "36",
                            "47",
                            "39",
                            "43"
                        ]
                    }
                ],
                "containerId": "20",
                "elements": [
                    {
                        "id": "33",
                        "x": 1925,
                        "y": 817
                    },
                    {
                        "id": "34",
                        "x": 1015,
                        "y": 817
                    },
                    {
                        "id": "4",
                        "x": 1925,
                        "y": 1307
                    },
                    {
                        "id": "17",
                        "x": 560,
                        "y": 10
                    },
                    {
                        "id": "6",
                        "x": 1015,
                        "y": 1307
                    },
                    {
                        "id": "18",
                        "x": 1470,
                        "y": 11
                    },
                    {
                        "id": "29",
                        "x": 105,
                        "y": 436
                    },
                    {
                        "id": "30",
                        "x": 1925,
                        "y": 436
                    },
                    {
                        "id": "31",
                        "x": 1015,
                        "y": 436
                    },
                    {
                        "id": "21",
                        "x": 105,
                        "y": 1307
                    },
                    {
                        "id": "32",
                        "x": 105,
                        "y": 817
                    }
                ],
                "relationships": [
                    {
                        "id": "40",
                        "position": 40
                    },
                    {
                        "id": "41",
                        "position": 55
                    },
                    {
                        "id": "42",
                        "position": 50
                    },
                    {
                        "id": "43"
                    },
                    {
                        "id": "37",
                        "position": 85
                    },
                    {
                        "id": "36",
                        "position": 45
                    },
                    {
                        "id": "35",
                        "position": 35
                    },
                    {
                        "id": "44"
                    },
                    {
                        "id": "45",
                        "position": 60
                    },
                    {
                        "id": "46"
                    },
                    {
                        "id": "47"
                    },
                    {
                        "id": "38",
                        "position": 85
                    },
                    {
                        "id": "39",
                        "position": 40
                    }
                ]
            }
        ],
        "dynamicViews": [
            {
                "description": "Summarises how the sign in feature works in the single-page application.",
                "key": "SignIn",
                "paperSize": "A5_Landscape",
                "elementId": "20",
                "relationships": [
                    {
                        "id": "35",
                        "description": "Submits credentials to",
                        "order": "1"
                    },
                    {
                        "id": "41",
                        "description": "Calls isAuthenticated() on",
                        "order": "2"
                    },
                    {
                        "id": "45",
                        "description": "select * from users where username = ?",
                        "order": "3"
                    }
                ],
                "elements": [
                    {
                        "id": "17",
                        "x": 552,
                        "y": 211
                    },
                    {
                        "id": "29",
                        "x": 1477,
                        "y": 211
                    },
                    {
                        "id": "32",
                        "x": 1477,
                        "y": 1116
                    },
                    {
                        "id": "21",
                        "x": 552,
                        "y": 1116
                    }
                ]
            }
        ],
        "deploymentViews": [
            {
                "softwareSystemId": "2",
                "description": "An example live deployment scenario for the Internet Banking System.",
                "key": "LiveDeployment",
                "paperSize": "A5_Landscape",
                "environment": "Live",
                "animations": [
                    {
                        "order": 1,
                        "elements": [
                            "66",
                            "67",
                            "65"
                        ]
                    },
                    {
                        "order": 2,
                        "elements": [
                            "63",
                            "64"
                        ]
                    },
                    {
                        "order": 3,
                        "elements": [
                            "68",
                            "69",
                            "70",
                            "71",
                            "73",
                            "74",
                            "75"
                        ],
                        "relationships": [
                            "77",
                            "72",
                            "76"
                        ]
                    },
                    {
                        "order": 4,
                        "elements": [
                            "78",
                            "79",
                            "80"
                        ],
                        "relationships": [
                            "81"
                        ]
                    },
                    {
                        "order": 5,
                        "elements": [
                            "82",
                            "83",
                            "84"
                        ],
                        "relationships": [
                            "85",
                            "86"
                        ]
                    }
                ],
                "elements": [
                    {
                        "id": "66",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "78",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "67",
                        "x": 150,
                        "y": 1026
                    },
                    {
                        "id": "79",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "68",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "69",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "80",
                        "x": 1820,
                        "y": 176
                    },
                    {
                        "id": "70",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "71",
                        "x": 985,
                        "y": 1026
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
                        "id": "84",
                        "x": 1820,
                        "y": 1026
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
                        "id": "63",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "75",
                        "x": 985,
                        "y": 176
                    },
                    {
                        "id": "64",
                        "x": 150,
                        "y": 176
                    },
                    {
                        "id": "65",
                        "x": 0,
                        "y": 0
                    }
                ],
                "relationships": [
                    {
                        "id": "72"
                    },
                    {
                        "id": "81"
                    },
                    {
                        "id": "86"
                    },
                    {
                        "id": "76"
                    },
                    {
                        "id": "85"
                    },
                    {
                        "id": "77"
                    }
                ]
            },
            {
                "softwareSystemId": "2",
                "description": "An example development deployment scenario for the Internet Banking System.",
                "key": "DevelopmentDeployment",
                "paperSize": "A5_Landscape",
                "environment": "Development",
                "animations": [
                    {
                        "order": 1,
                        "elements": [
                            "59",
                            "60",
                            "50"
                        ]
                    },
                    {
                        "order": 2,
                        "elements": [
                            "51",
                            "52",
                            "53",
                            "54"
                        ],
                        "relationships": [
                            "61",
                            "62"
                        ]
                    },
                    {
                        "order": 3,
                        "elements": [
                            "55",
                            "56",
                            "57"
                        ],
                        "relationships": [
                            "58"
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
                        "id": "56",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "57",
                        "x": 1840,
                        "y": 834
                    },
                    {
                        "id": "59",
                        "x": 0,
                        "y": 0
                    },
                    {
                        "id": "60",
                        "x": 140,
                        "y": 664
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
                        "id": "53",
                        "x": 990,
                        "y": 494
                    },
                    {
                        "id": "54",
                        "x": 990,
                        "y": 834
                    }
                ],
                "relationships": [
                    {
                        "id": "61"
                    },
                    {
                        "id": "62"
                    },
                    {
                        "id": "58",
                        "position": 50
                    }
                ]
            }
        ],
        "configuration": {
            "branding": {},
            "styles": {
                "elements": [
                    {
                        "tag": "Element"
                    },
                    {
                        "tag": "Software System",
                        "background": "#1168bd",
                        "color": "#ffffff"
                    },
                    {
                        "tag": "Container",
                        "background": "#438dd5",
                        "color": "#ffffff"
                    },
                    {
                        "tag": "Component",
                        "background": "#85bbf0",
                        "color": "#000000"
                    },
                    {
                        "tag": "Person",
                        "background": "#08427b",
                        "color": "#ffffff",
                        "fontSize": 22,
                        "shape": "Person"
                    },
                    {
                        "tag": "Existing System",
                        "background": "#999999",
                        "color": "#ffffff"
                    },
                    {
                        "tag": "Bank Staff",
                        "background": "#999999",
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
                        "tag": "Failover",
                        "opacity": 25
                    }
                ],
                "relationships": [
                    {
                        "tag": "Failover",
                        "position": 70,
                        "opacity": 25
                    }
                ]
            },
            "terminology": {},
            "lastSavedView": "Components",
            "themes": []
        },
        "filteredViews": []
    }
}
`
