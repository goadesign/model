package model

import (
	. "goa.design/structurizr/dsl"
	"goa.design/structurizr/expr"
)

var _ = Workspace("Big Bank plc", "This is an example workspace to illustrate the key features of Structurizr, based around a fictional online banking system.", func() {
	Enterprise("Big Bank plc")

	var BackOfficeStaff = Person("Back Office Staff", "Administration and support staff within the bank.", func() {
		Uses("Mainframe Banking System", "Uses", Synchronous, func() {
			Tag("Relationship", "Synchronous")
		})
		Tag("Element", "Person", "Bank Staff")
	})

	var CustomerServiceStaff = Person("Customer Service Staff", "Customer service staff within the bank.", func() {
		Uses("Mainframe Banking System", "Uses", Synchronous, func() {
			Tag("Relationship", "Synchronous")
		})
		Tag("Element", "Person", "Bank Staff")
	})

	var PersonalBankingCustomer = Person("Personal Banking Customer", "A customer of the bank, with personal bank accounts.", func() {
		External()
		InteractsWith("Customer Service Staff", "Asks questions to", "Telephone", Synchronous, func() {
			Tag("Relationship", "Synchronous")
		})
		Uses("Single Page Application", "Views account balances, and makes payments using", Synchronous, func() {
			Tag("Relationship", "Synchronous")
		})
		Uses("ATM", "Withdraws cash using", Synchronous, func() {
			Tag("Relationship", "Synchronous")
		})
		Uses("Internet Banking System", "Views account balances, and makes payments using", Synchronous, func() {
			Tag("Relationship", "Synchronous")
		})
		Uses("Mobile Application", "Views account balances, and makes payments using", Synchronous, func() {
			Tag("Relationship", "Synchronous")
		})
		Uses("Web Application", "Visits bigbank.com/ib using", "HTTPS", Synchronous, func() {
			Tag("Relationship", "Synchronous")
		})
		Tag("Element", "Person")
	})

	var MainframeBankingSystem = SoftwareSystem("Mainframe Banking System", "Stores all of the core banking information about customers, accounts, transactions, etc.", func() {
		Tag("Element", "Software System", "Existing System")
	})

	var ATM = SoftwareSystem("ATM", "Allows customers to withdraw cash.", func() {
		Uses(MainframeBankingSystem, "Uses", Synchronous, func() {
			Tag("Relationship", "Synchronous")
		})
		Tag("Element", "Software System", "Existing System")
	})

	var EMailSystem = SoftwareSystem("E-mail System", "The internal Microsoft Exchange e-mail system.", func() {
		Delivers(PersonalBankingCustomer, "Sends e-mails to", Synchronous, func() {
			Tag("Relationship", "Synchronous")
		})
		Tag("Element", "Software System", "Existing System")
	})

	var (
		// Forward declaration so variables can be used to define views.
		Database       *expr.Container
		APIApplication *expr.Container
		MobileApp      *expr.Container
		SinglePageApp  *expr.Container
		WebApp         *expr.Container
	)

	var InternetBankingSystem = SoftwareSystem("Internet Banking System", "Allows customers to view information about their bank accounts, and make payments.", func() {
		Uses("Email System", "Sends e-mail using", Synchronous, func() {
			Tag("Relationship", "Synchronous")
		})
		Uses("Mainframe Banking System", "Gets account information from, and makes payments using", "", Synchronous, func() {
			Tag("Relationship", "Synchronous")
		})
		Tag("Element", "Software System")

		Database = Container("Database", "Stores user registration information, hashed authentication credentials, access logs, etc.", "Oracle Database Schema", func() {
			Tag("Element", "Container", "Database")
		})

		APIApplication = Container("API Application", "Provides Internet banking functionality via a JSON/HTTPS API.", "Java and Spring MVC", func() {
			Uses(MainframeBankingSystem, "Makes API calls to", "XML/HTTPS", Synchronous, func() {
				Tag("Relationship", "Synchronous")
			})
			Uses(Database, "Reads from and writes to", "JDBC", Synchronous, func() {
				Tag("Relationship", "Synchronous")
			})
			Uses(EMailSystem, "Sends e-mail using", "SMTP", Synchronous, func() {
				Tag("Relationship", "Synchronous")
			})
			Tag("Element", "Container")
			Component("Mainframe Banking System Facade", "A facade onto the mainframe banking system.", "Spring Bean", func() {
				Uses(MainframeBankingSystem, "Uses", "XML/HTTPS", Synchronous, func() {
					Tag("Relationship", "Synchronous")
				})
				Tag("Element", "Component")
			})
			Component("Accounts Summary Controller", "Provides customers with a summary of their bank accounts.", "Spring MVC Rest Controller", func() {
				Uses("Mainframe Banking System Facade", "Uses", Synchronous, func() {
					Tag("Relationship", "Synchronous")
				})
				Tag("Element", "Component")
			})
			Component("Email Component", "Sends e-mails to users.", "Spring Bean", func() {
				Uses(EMailSystem, "Sends e-mail using", Synchronous, func() {
					Tag("Relationship", "Synchronous")
				})
				Tag("Element", "Component")
			})
			Component("Security Component", "Provides functionality related to signing in, changing passwords, etc..", "Spring Bean", func() {
				Uses(Database, "Reads from and writes to", "JDBC", Synchronous, func() {
					Tag("Relationship", "Synchronous")
				})
				Tag("Element", "Component")
			})
			Component("Reset Password Controller", "Allows users to reset their passwords with a single use URL.", "Spring MVC Rest Controller", func() {
				Uses("Email Component", "Uses", Synchronous, func() {
					Tag("Relationship", "Synchronous")
				})
				Uses("Security Component", "Uses", Synchronous, func() {
					Tag("Relationship", "Synchronous")
				})
				Tag("Element", "Component")
			})
			Component("Sign In Controller", "Allows users to sign in to the Internet Banking System.", "Spring MVC Rest Controller", func() {
				Uses("Security Component", "Uses", Synchronous, func() {
					Tag("Relationship", "Synchronous")
				})
				Tag("Element", "Component")
			})
		})

		MobileApp = Container("Mobile App", "Provides a limited subset of the Internet banking functionality to customers via their mobile device.", "Xamarin", func() {
			Uses("Reset Password Controller", "Makes API calls to", "JSON/HTTPS", Synchronous, func() {
				Tag("Relationship", "Synchronous")
			})
			Uses(APIApplication, "Makes API calls to", "JSON/HTTPS", Synchronous, func() {
				Tag("Relationship", "Synchronous")
			})
			Uses("Sign In Controller", "Makes API calls to", "JSON/HTTPS", Synchronous, func() {
				Tag("Relationship", "Synchronous")
			})
			Uses("Accounts Summary Controller", "Makes API calls to", "JSON/HTTPS", Synchronous, func() {
				Tag("Relationship", "Synchronous")
			})
			Tag("Element", "Container", "Mobile App")
		})

		SinglePageApp = Container("Single-Page Application", "Provides all of the Internet banking functionality to customers via their web browser.", "JavaScript and Angular", func() {
			Uses("Accounts Summary Controller", "Makes API calls to", "JSON/HTTPS", Synchronous, func() {
				Tag("Relationship", "Synchronous")
			})
			Uses("Sign In Controller", "Makes API calls to", "JSON/HTTPS", Synchronous, func() {
				Tag("Relationship", "Synchronous")
			})
			Uses(APIApplication, "Makes API calls to", "JSON/HTTPS", Synchronous, func() {
				Tag("Relationship", "Synchronous")
			})
			Uses("Reset Password Controller", "Makes API calls to", "JSON/HTTPS", Synchronous, func() {
				Tag("Relationship", "Synchronous")
			})
		})

		WebApp = Container("Web Application", "Delivers the static content and the Internet banking single page application.", "Java and Spring MVC", func() {
			Uses(SinglePageApp, "Delivers to the customer's web browser", Synchronous, func() {
				Tag("Relationship", "Synchronous")
			})
		})
	})

	var (
		// Forward declarations so variables can be used in deployment views.
		DevLaptop          *expr.DeploymentNode
		DevDBDocker        *expr.DeploymentNode
		DevDB              *expr.DeploymentNode
		DevDBInstance      *expr.ContainerInstance
		DevWebServerDocker *expr.DeploymentNode
		DevWebServer       *expr.DeploymentNode
		DevAPIAppInstance  *expr.ContainerInstance
		DevWebAppInstance  *expr.ContainerInstance
		DevWebBrowser      *expr.DeploymentNode
		DevSPAInstance     *expr.ContainerInstance
	)

	DeploymentEnvironment("Development", func() {

		DevLaptop = DeploymentNode("Developer Laptop", "A developer laptop", "Microsoft Windows 10 or Apple macOS", func() {
			Tag("Element", "Deployment Node")

			DevDBDocker = DeploymentNode("Docker Container - Database Server", "A Docker container.", "Docker", func() {
				Tag("Element", "Deployment Node")

				DevDB = DeploymentNode("Database Server", "A development database.", "Oracle 12c", func() {
					DevDBInstance = ContainerInstance(Database, func() {
						Tag("Container Instance")
					})
					Tag("Element", "Deployment Node")
				})
			})

			DevWebServerDocker = DeploymentNode("Docker Container - Web Server", "A Docker container.", "Docker", func() {
				Tag("Element", "Deployment Node")

				DevWebServer = DeploymentNode("Apache Tomcat", "An open source Java EE web server.", "Apache Tomcat 8.x", func() {
					DevAPIAppInstance = ContainerInstance(APIApplication, func() {
						Tag("Container Instance")
					})
					DevWebAppInstance = ContainerInstance(WebApp, func() {
						Tag("Container Instance")
					})
					Tag("Element", "Deployment Node")
					Prop("Java Version", "8")
					Prop("Xms", "1024M")
					Prop("Xmx", "512M")
				})
			})

			DevWebBrowser = DeploymentNode("Web Browser", "", "Chrome, Firefox, Safari or Edge", func() {
				Tag("Element", "Deployment Node")

				DevSPAInstance = ContainerInstance(SinglePageApp, func() {
					Tag("Container Instance")
				})
			})
		})
	})

	var (
		// Forward declarations so variables can be used in deployment views.
		DataCenter          *expr.DeploymentNode
		BigBankAPI          *expr.DeploymentNode
		LiveAPIApp          *expr.DeploymentNode
		LiveAPIAppInstance  *expr.ContainerInstance
		BigBankDB01         *expr.DeploymentNode
		PrimaryDB           *expr.DeploymentNode
		PrimaryDBInstance   *expr.ContainerInstance
		BigBankDB02         *expr.DeploymentNode
		SecondaryDB         *expr.DeploymentNode
		SecondaryDBInstance *expr.ContainerInstance
		BigBankWeb          *expr.DeploymentNode
		LiveWebApp          *expr.DeploymentNode
		LiveWebAppInstance  *expr.ContainerInstance
		CustomerComputer    *expr.DeploymentNode
		LiveWebBrowser      *expr.DeploymentNode
		CustomerSPA         *expr.ContainerInstance
		CustomerMobile      *expr.DeploymentNode
		CustomerMobileApp   *expr.ContainerInstance
	)

	DeploymentEnvironment("Live", func() {

		DataCenter = DeploymentNode("Big Bank plc", "", "Big Bank plc data center", func() {
			Tag("Element", "Deployment Node")

			BigBankAPI = DeploymentNode("bigbank-api***", "A web server residing in the web server farm, accessed via F5 BIG-IP LTMs.", "Ubuntu 16.04 LTS", func() {
				Tag("Element", "Deployment Node")
				Instances(8)
				Prop("Location", "London and Reading")

				LiveAPIApp = DeploymentNode("Apache Tomcat", "An open source Java EE web server.", "Apache Tomcat 8.x", func() {
					Tag("Element", "Deployment Node")
					Prop("Java Version", "8")
					Prop("Xms", "1024M")
					Prop("Xmx", "512M")

					LiveAPIAppInstance = ContainerInstance(APIApplication, func() {
						Tag("Container Instance")
						InstanceID(2)
					})
				})
			})

			BigBankDB01 = DeploymentNode("bigbank-db01", "The primary database server.", "Ubuntu 16.04 LTS", func() {
				Tag("Element", "Deployment Node")
				Prop("Location", "London")

				PrimaryDB = DeploymentNode("Oracle - Primary", "The primary, live database server.", "Oracle 12c", func() {
					Tag("Element", "Deployment Node")

					PrimaryDBInstance = ContainerInstance(Database, func() {
						Tag("Container Instance")
						InstanceID(2)
					})
				})
			})

			BigBankDB02 = DeploymentNode("bigbank-db02", "The secondary database server.", "Ubuntu 16.04 LTS", func() {
				Tag("Element", "Deployment Node")
				Prop("Location", "Reading")

				SecondaryDB = DeploymentNode("Oracle - Secondary", "A secondary, standby database server, used for failover purposes only.", "Oracle 12c", func() {
					Tag("Element", "Deployment Node", "Failover")

					SecondaryDBInstance = ContainerInstance(Database, func() {
						Tag("Container Instance", "Failover")
						InstanceID(3)
					})
				})
			})

			BigBankWeb = DeploymentNode("bigbank-web***", "A web server residing in the web server farm, accessed via F5 BIG-IP LTMs.", "Ubuntu 16.04 LTS", func() {
				Tag("Element", "Deployment Node")
				Instances(4)
				Prop("Location", "London and Reading")

				LiveWebApp = DeploymentNode("Apache Tomcat", "An open source Java EE web server.", "Apache Tomcat 8.x", func() {
					Tag("Element", "Deployment Node")
					Prop("Java Version", "8")
					Prop("Xms", "1024M")
					Prop("Xmx", "512M")

					LiveWebAppInstance = ContainerInstance(WebApp, func() {
						Tag("Container Instance")
						InstanceID(2)
					})
				})
			})
		})

		CustomerComputer = DeploymentNode("Customer's computer", "", "Microsoft Windows or Apple macOS", func() {
			Tag("Element", "Deployment Node")

			LiveWebBrowser = DeploymentNode("Web Browser", "", "Chrome, Firefox, Safari or Edge", func() {
				Tag("Element", "Deployment Node")

				CustomerSPA = ContainerInstance(SinglePageApp, func() {
					Tag("Container Instance")
				})
			})
		})

		CustomerMobile = DeploymentNode("Customer's mobile device", "", "Apple iOS or Android", func() {
			Tag("Element", "Deployment Node")

			CustomerMobileApp = ContainerInstance(MobileApp, func() {
				Tag("Container Instance")
			})
		})
	})

	Views(func() {

		SystemLandscapeView("SystemLandscape", "The system landscape diagram for Big Bank plc.", func() {
			PaperSize(SizeA5Landscape)
			Animation("Personal Banking Customer", "Internet Banking System", "Mainframe Banking System", "Email System")
			Animation("ATM")
			Animation("Customer Service Staff", "Bank Office Staff")
			EnterpriseBoundaryVisible()

			Add(PersonalBankingCustomer, func() {
				Coord(87, 643)
			})
			Add(InternetBankingSystem, func() {
				Coord(1012, 813)
			})
			Add(MainframeBankingSystem, func() {
				Coord(1922, 693)
			})
			Add(EMailSystem, func() {
				Coord(1012, 1326)
			})
			Add(ATM, func() {
				Coord(1012, 301)
			})
			Add(CustomerServiceStaff, func() {
				Coord(1947, 36)
			})
			Add(BackOfficeStaff, func() {
				Coord(1947, 1241)
			})

			Link(PersonalBankingCustomer, CustomerServiceStaff, func() {
				Vertices(285, 240)
			})
		})

		ContainerView("Internet Banking System", "Containers", "The container diagram for the Internet Banking System.", func() {
			PaperSize(SizeA5Landscape)
			Animation("Personal Banking Customer", "Mainframe Banking System", "Email System")
			Animation("Web Application")
			Animation("Single-Page Application")
			Animation("Mobile App")
			Animation("API Application")
			Animation("Database")

			Add(PersonalBankingCustomer, func() {
				Coord(1056, 24)
			})
			Add(MainframeBankingSystem, func() {
				Coord(2012, 1214)
			})
			Add(SinglePageApp, func() {
				Coord(780, 664)
			})
			Add(EMailSystem, func() {
				Coord(2012, 664)
			})
			Add(MobileApp, func() {
				Coord(1283, 664)
			})
			Add(WebApp, func() {
				Coord(37, 664)
			})
			Add(APIApplication, func() {
				Coord(1031, 1214)
			})
			Add(Database, func() {
				Coord(37, 1214)
			})
		})

		ComponentView("API Application", "Components", "The component diagram for the API Application", func() {
			PaperSize(SizeA5Landscape)
			Animation("Mainframe Banking System", "Single-Page Application", "Email System", "Mobile App", "Database")
			Animation("Sign In Controller", "Security Component")
			Animation("Mainframe Banking System Facade", "Accounts Summary Controller")
			Animation("Email Component", "Reset Password Controller")

			Add("MainframeBankingSystemFacade", func() {
				Coord(1925, 817)
			})
			Add("Email Component", func() {
				Coord(1015, 817)
			})
			Add(MainframeBankingSystem, func() {
				Coord(1925, 1307)
			})
			Add(SinglePageApp, func() {
				Coord(560, 10)
			})
			Add(EMailSystem, func() {
				Coord(1015, 1307)
			})
			Add(MobileApp, func() {
				Coord(1470, 11)
			})
			Add("Sign In Controller", func() {
				Coord(105, 436)
			})
			Add("Accounts Summary Controller", func() {
				Coord(1925, 436)
			})
			Add("Reset Password Controller", func() {
				Coord(1015, 436)
			})
			Add(Database, func() {
				Coord(105, 1307)
			})
			Add("Security Component", func() {
				Coord(105, 817)
			})
		})

		DynamicView("API Application", "SignIn", "Summarises how the sign in feature works in the single-page application.", func() {
			PaperSize(SizeA5Landscape)
			Link(SinglePageApp, "SignInController", func() {
				Description("Submits credentials to")
			})
			Link("Sign In Controller", "Security Component", func() {
				Description("Calls isAuthenticated() on")
			})
			Link("Security Component", Database, func() {
				Description("select * from users where username = ?")
			})
		})

		DeploymentView("Internet Banking System", "Live", "LiveDeployment", "An example live deployment scenario for the Internet Banking System.", func() {
			PaperSize(SizeA5Landscape)
			AutoLayout(RankBottomTop)

			Animation(LiveWebBrowser, CustomerSPA, CustomerComputer)
			Animation(CustomerMobile, CustomerMobileApp)
			Animation(DataCenter, BigBankWeb, LiveWebApp, LiveWebAppInstance, BigBankAPI, LiveAPIApp, LiveAPIAppInstance)
			Animation(BigBankDB01, PrimaryDB, PrimaryDBInstance)
			Animation(BigBankDB02, SecondaryDB, SecondaryDBInstance)

			Add(LiveWebBrowser)
			Add(CustomerSPA)
			Add(CustomerComputer)
			Add(CustomerMobile)
			Add(CustomerMobileApp)
			Add(DataCenter)
			Add(BigBankWeb)
			Add(LiveWebApp)
			Add(LiveWebAppInstance)
			Add(BigBankAPI)
			Add(LiveAPIApp)
			Add(LiveAPIAppInstance)
			Add(BigBankDB01)
			Add(PrimaryDB)
			Add(PrimaryDBInstance)
			Add(BigBankDB02)
			Add(SecondaryDB)
			Add(SecondaryDBInstance)
		})

		DeploymentView("Internet Banking System", "Development", "DevelopmentDeployment", "An example development deployment scenario for the Internet Banking System.", func() {
			PaperSize(SizeA5Landscape)
			AutoLayout(RankBottomTop)

			Animation(DevWebBrowser, DevSPAInstance, DevLaptop)
			Animation(DevWebServerDocker, DevWebServer, DevAPIAppInstance, DevWebAppInstance)
			Animation(DevDBDocker, DevDB, DevDBInstance)

			Add(DevWebBrowser)
			Add(DevSPAInstance)
			Add(DevLaptop)
			Add(DevWebServerDocker)
			Add(DevWebServer)
			Add(DevAPIAppInstance)
			Add(DevWebAppInstance)
			Add(DevDBDocker)
			Add(DevDB)
			Add(DevDBInstance)
		})
	})

	Styles(func() {
		ElementStyle("Software System", func() {
			Background("#1168bd")
			Color("#ffffff")
		})
		ElementStyle("Container", func() {
			Background("#438dd5")
			Color("#ffffff")
		})
		ElementStyle("Component", func() {
			Background("#85bbf0")
			Color("#000000")
		})
		ElementStyle("Person", func() {
			Background("#08427b")
			Color("#ffffff")
			FontSize(22)
			Shape(ShapePerson)
		})
		ElementStyle("Existing System", func() {
			Background("#999999")
			Color("#ffffff")
		})
		ElementStyle("Bank Staff", func() {
			Background("#999999")
			Color("#ffffff")
		})
		ElementStyle("Web Browser", func() {
			Shape(ShapeWebBrowser)
		})
		ElementStyle("Mobile App", func() {
			Shape(expr.ShapeMobileDeviceLandscape)
		})
		ElementStyle("Database", func() {
			Shape(ShapeCylinder)
		})
		ElementStyle("Failover", func() {
			Opacity(25)
		})
		RelationshipStyle("Failover", func() {
			Position(70)
			Opacity(25)
		})
	})
})
