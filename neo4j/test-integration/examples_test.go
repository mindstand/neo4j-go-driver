/*
 * Copyright (c) "Neo4j"
 * Neo4j Sweden AB [http://neo4j.com]
 *
 * This file is part of Neo4j.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package test_integration

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/test-integration/dbserver"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Examples", func() {

	Context("Single Instance", func() {
		var (
			uri      string
			username string
			password string
		)

		server := dbserver.GetDbServer()

		BeforeEach(func() {
			uri = server.BoltURI()
			username = server.Username
			password = server.Password
		})

		Specify("Hello World", func() {
			greeting, err := helloWorld(uri, username, password)

			Expect(err).To(BeNil())
			Expect(greeting).To(ContainSubstring("hello, world"))
		})

		Specify("Driver Lifecycle", func() {
			driver, err := createDriver(uri, username, password)
			Expect(err).To(BeNil())
			Expect(driver).NotTo(BeNil())

			err = closeDriver(driver)
			Expect(err).To(BeNil())
		})

		Specify("Basic Authentication", func() {
			driver, err := createDriverWithBasicAuth(uri, username, password)
			Expect(err).To(BeNil())
			Expect(driver).NotTo(BeNil())

			err = driver.Close()
			Expect(err).To(BeNil())
		})

		Specify("Config - With Max Retry Time", func() {
			driver, err := createDriverWithMaxRetryTime(uri, username, password)
			Expect(err).To(BeNil())
			Expect(driver).NotTo(BeNil())

			err = driver.Close()
			Expect(err).To(BeNil())
		})

		Specify("Config - With Customized Connection Pool", func() {
			driver, err := createDriverWithCustomizedConnectionPool(uri, username, password)
			Expect(err).To(BeNil())
			Expect(driver).NotTo(BeNil())

			err = driver.Close()
			Expect(err).To(BeNil())
		})

		Specify("Config - With Connection Timeout", func() {
			driver, err := createDriverWithConnectionTimeout(uri, username, password)
			Expect(err).To(BeNil())
			Expect(driver).NotTo(BeNil())

			err = driver.Close()
			Expect(err).To(BeNil())
		})

		/*
			Specify("Service Unavailable", func() {
				driver, err := createDriverWithMaxRetryTime("bolt://localhost:8080", username, password)
				Expect(err).To(BeNil())
				Expect(driver).NotTo(BeNil())
				defer driver.Close()

				err = createItem(driver)
				errDescr := err.Error()
				Expect(errDescr).To(ContainSubstring("retryable operation failed to complete after"))
				//Expect(err).To(test.BeGenericError(ContainSubstring("retryable operation failed to complete after")))
			})
		*/

		Specify("Session", func() {
			driver, err := createDriverWithMaxRetryTime(uri, username, password)
			Expect(err).To(BeNil())
			Expect(driver).NotTo(BeNil())
			defer driver.Close()

			err = addPersonInSession(driver, "Tom")
			Expect(err).To(BeNil())
			count, err := countNodes(driver, "Person", "name", "Tom")
			Expect(err).To(BeNil())
			Expect(count).To(BeNumerically("==", 1))
		})

		Specify("Autocommit Transaction", func() {
			driver, err := createDriverWithMaxRetryTime(uri, username, password)
			Expect(err).To(BeNil())
			Expect(driver).NotTo(BeNil())
			defer driver.Close()

			err = addPersonInAutoCommitTx(driver, "Shanon")
			Expect(err).To(BeNil())
			count, err := countNodes(driver, "Person", "name", "Shanon")
			Expect(err).To(BeNil())
			Expect(count).To(BeNumerically("==", 1))
		})

		Specify("Pass Bookmarks", func() {
			driver, err := createDriverWithMaxRetryTime(uri, username, password)
			Expect(err).To(BeNil())
			Expect(driver).NotTo(BeNil())
			defer driver.Close()

			err = addEmployAndMakeFriends(driver)
			Expect(err).To(BeNil())

			count, err := countNodes(driver, "Person", "name", "Alice")
			Expect(err).To(BeNil())
			Expect(count).To(BeNumerically("==", 1))

			count, err = countNodes(driver, "Person", "name", "Bob")
			Expect(err).To(BeNil())
			Expect(count).To(BeNumerically("==", 1))

			count, err = countNodes(driver, "Company", "name", "LexCorp")
			Expect(err).To(BeNil())
			Expect(count).To(BeNumerically("==", 1))

			count, err = countNodes(driver, "Company", "name", "Wayne Enterprises")
			Expect(err).To(BeNil())
			Expect(count).To(BeNumerically("==", 1))
		})

		Specify("Read/Write Transaction", func() {
			driver, err := createDriverWithMaxRetryTime(uri, username, password)
			Expect(err).To(BeNil())
			Expect(driver).NotTo(BeNil())
			defer driver.Close()

			id, err := addPersonNode(driver, "Jason")
			Expect(err).To(BeNil())
			Expect(id).To(BeNumerically(">=", 0))
		})

		Specify("Get People", func() {
			driver, err := createDriverWithMaxRetryTime(uri, username, password)
			Expect(err).To(BeNil())
			Expect(driver).NotTo(BeNil())
			defer driver.Close()

			id, err := addPersonNode(driver, "Annie")
			Expect(err).To(BeNil())
			Expect(id).To(BeNumerically(">=", 0))

			id, err = addPersonNode(driver, "Joe")
			Expect(err).To(BeNil())
			Expect(id).To(BeNumerically(">=", 0))

			people, err := getPeople(driver)
			Expect(err).To(BeNil())
			Expect(people).To(ContainElement("Annie"))
			Expect(people).To(ContainElement("Joe"))
		})

		Specify("Result Retain", func() {
			driver, err := createDriverWithMaxRetryTime(uri, username, password)
			Expect(err).To(BeNil())
			Expect(driver).NotTo(BeNil())
			defer driver.Close()

			id, err := addPersonNode(driver, "Carl")
			Expect(err).To(BeNil())
			Expect(id).To(BeNumerically(">=", 0))

			id, err = addPersonNode(driver, "Thomas")
			Expect(err).To(BeNil())
			Expect(id).To(BeNumerically(">=", 0))

			count, err := addPersonsAsEmployees(driver, "Acme")
			Expect(err).To(BeNil())
			Expect(count).To(BeNumerically(">=", 2))
		})

		Specify("Point2D", func() {
			driver, err := createDriverWithMaxRetryTime(uri, username, password)
			Expect(err).To(BeNil())
			Expect(driver).NotTo(BeNil())
			defer driver.Close()

			// tag::geospatial-types-point2d[]
			// Creating a 2D point in Cartesian space
			cartesian := dbtype.Point2D{
				X:            2.5,
				Y:            -2,
				SpatialRefId: 7203,
			}

			// Creating a 2D point in WGS84 space
			wgs84 := dbtype.Point2D{
				X:            -1.5,
				Y:            1,
				SpatialRefId: 4326,
			}
			// end::geospatial-types-point2d[]

			session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
			Expect(session).ToNot(BeNil())
			defer session.Close()

			recordWithCartesian, err := echo(session, cartesian)
			Expect(err).To(BeNil())
			Expect(recordWithCartesian).ToNot(BeNil())

			recordWithWgs84, err := echo(session, wgs84)
			Expect(err).To(BeNil())
			Expect(recordWithWgs84).ToNot(BeNil())

			// tag::geospatial-types-point2d[]

			// Reading a Cartesian point from a record
			field, _ := recordWithCartesian.Get("fieldName")
			fieldCartesian, _ := field.(dbtype.Point2D)

			// Serializing
			fieldCartesian.String() // Point{srId=7203, x=2.500000, y=-2.000000}

			// Acessing members
			print(fieldCartesian.X)            // 2.500000
			print(fieldCartesian.Y)            // -2.000000
			print(fieldCartesian.SpatialRefId) // 7203

			// Reading a WGS84 point from a record
			field, _ = recordWithWgs84.Get("fieldName")
			fieldWgs84 := field.(dbtype.Point2D)

			// Serializing
			fieldWgs84.String() // Point{srId=4326, x=-1.500000, y=1.00000}

			// Acessing members
			print(fieldWgs84.X)            // -1.500000
			print(fieldWgs84.Y)            // 1.000000
			print(fieldWgs84.SpatialRefId) // 4326
			// end::geospatial-types-point2d[]

			Expect(fieldCartesian.String()).To(Equal("Point{srId=7203, x=2.500000, y=-2.000000}"))
			Expect(fieldCartesian.X).To(Equal(cartesian.X))
			Expect(fieldCartesian.Y).To(Equal(cartesian.Y))
			Expect(fieldCartesian.SpatialRefId).To(Equal(cartesian.SpatialRefId))

			Expect(fieldWgs84.String()).To(Equal("Point{srId=4326, x=-1.500000, y=1.000000}"))
			Expect(fieldWgs84.X).To(Equal(wgs84.X))
			Expect(fieldWgs84.Y).To(Equal(wgs84.Y))
			Expect(fieldWgs84.SpatialRefId).To(Equal(wgs84.SpatialRefId))
		})

		Specify("Point3D", func() {
			driver, err := createDriverWithMaxRetryTime(uri, username, password)
			Expect(err).To(BeNil())
			Expect(driver).NotTo(BeNil())
			defer driver.Close()

			// tag::geospatial-types-point3d[]
			// Creating a 3D point in Cartesian space
			cartesian := dbtype.Point3D{
				X:            2.5,
				Y:            -2,
				Z:            2,
				SpatialRefId: 9157,
			}

			// Creating a 3D point in WGS84 space
			wgs84 := dbtype.Point3D{
				X:            -1.5,
				Y:            1,
				Z:            3,
				SpatialRefId: 4979,
			}
			// end::geospatial-types-point3d[]

			session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
			Expect(session).ToNot(BeNil())
			defer session.Close()

			recordWithCartesian, err := echo(session, cartesian)
			Expect(err).To(BeNil())
			Expect(recordWithCartesian).ToNot(BeNil())

			recordWithWgs84, err := echo(session, wgs84)
			Expect(err).To(BeNil())
			Expect(recordWithWgs84).ToNot(BeNil())

			// tag::geospatial-types-point3d[]

			// Reading a Cartesian point from a record
			field, _ := recordWithCartesian.Get("fieldName")
			fieldCartesian := field.(dbtype.Point3D)

			// Serializing
			fieldCartesian.String() // Point{srId=9157, x=2.500000, y=-2.000000, z=2.000000}

			// Acessing members
			print(fieldCartesian.X)            // 2.500000
			print(fieldCartesian.Y)            // -2.000000
			print(fieldCartesian.Z)            // 2.000000
			print(fieldCartesian.SpatialRefId) // 7203

			// Reading a WGS84 point from a record
			field, _ = recordWithWgs84.Get("fieldName")
			fieldWgs84 := field.(dbtype.Point3D)

			// Serializing
			fieldWgs84.String() // Point{srId=4979, x=-1.500000, y=1.00000, z=3.000000}

			// Acessing members
			print(fieldWgs84.X)            // -1.500000
			print(fieldWgs84.Y)            // 1.000000
			print(fieldWgs84.Z)            // 3.000000
			print(fieldWgs84.SpatialRefId) // 4979
			// end::geospatial-types-point3d[]

			Expect(fieldCartesian.String()).To(Equal("Point{srId=9157, x=2.500000, y=-2.000000, z=2.000000}"))
			Expect(fieldCartesian.X).To(Equal(cartesian.X))
			Expect(fieldCartesian.Y).To(Equal(cartesian.Y))
			Expect(fieldCartesian.Z).To(Equal(cartesian.Z))
			Expect(fieldCartesian.SpatialRefId).To(Equal(cartesian.SpatialRefId))

			Expect(fieldWgs84.String()).To(Equal("Point{srId=4979, x=-1.500000, y=1.000000, z=3.000000}"))
			Expect(fieldWgs84.X).To(Equal(wgs84.X))
			Expect(fieldWgs84.Y).To(Equal(wgs84.Y))
			Expect(fieldWgs84.Z).To(Equal(wgs84.Z))
			Expect(fieldWgs84.SpatialRefId).To(Equal(wgs84.SpatialRefId))
		})
	})

	/*
		Context("Causal Cluster", func() {
			server := dbserver.GetDbServer()
			username := server.Username
			password := server.Password

			Specify("Config - Address Resolver", func() {
				if !server.IsCluster {
					Skip("Need cluster")
				}
				var addresses []neo4j.ServerAddress
				for _, server := range server.ClusterMembers {
					addresses = append(addresses, &url.URL{Host: server.HostnameAndPort})
				}

				driver, err := createDriverWithAddressResolver("neo4j://x.acme.com", username, password, addresses...)
				Expect(err).To(BeNil())
				Expect(driver).NotTo(BeNil())

				err = createItem(driver)
				Expect(err).To(BeNil())

				err = driver.Close()
				Expect(err).To(BeNil())
			})
		})
	*/

})

// tag::hello-world[]
func helloWorld(uri, username, password string) (string, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return "", err
	}
	defer driver.Close()

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	greeting, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"CREATE (a:Greeting) SET a.message = $message RETURN a.message + ', from node ' + id(a)",
			map[string]interface{}{"message": "hello, world"})
		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return "", err
	}

	return greeting.(string), nil
}

// end::hello-world[]

// tag::driver-lifecycle[]
func createDriver(uri, username, password string) (neo4j.Driver, error) {
	return neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
}

// call on application exit
func closeDriver(driver neo4j.Driver) error {
	return driver.Close()
}

// end::driver-lifecycle[]

// tag::basic-auth[]
func createDriverWithBasicAuth(uri, username, password string) (neo4j.Driver, error) {
	return neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
}

// end::basic-auth[]

// tag::kerberos-auth[]
func createDriverWithKerberosAuth(uri, ticket string) (neo4j.Driver, error) {
	return neo4j.NewDriver(uri, neo4j.KerberosAuth(ticket))
}

// end::kerberos-auth[]

// tag::custom-auth[]
func createDriverWithCustomAuth(uri, principal, credentials, realm, scheme string, parameters map[string]interface{}) (neo4j.Driver, error) {
	return neo4j.NewDriver(uri, neo4j.CustomAuth(scheme, principal, credentials, realm, parameters))
}

// end::custom-auth[]

// tag::config-unencrypted[]
/*
func createDriverWithoutEncryption(uri, username, password string) (neo4j.Driver, error) {
	return neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""), func(config *neo4j.Config) {
	})
}

// end::config-unencrypted[]

// tag::config-trust[]
func createDriverWithTrustStrategy(uri, username, password string) (neo4j.Driver, error) {
	return neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""), func(config *neo4j.Config) {
		config.TrustStrategy = neo4j.TrustAny(true)
	})
}

// end::config-trust[]
*/

// tag::config-custom-resolver[]
func createDriverWithAddressResolver(virtualURI, username, password string, addresses ...neo4j.ServerAddress) (neo4j.Driver, error) {
	// Address resolver is only valid for neo4j uri
	return neo4j.NewDriver(virtualURI, neo4j.BasicAuth(username, password, ""), func(config *neo4j.Config) {
		config.AddressResolver = func(address neo4j.ServerAddress) []neo4j.ServerAddress {
			return addresses
		}
	})
}

func addPerson(name string) error {
	const (
		username = "neo4j"
		password = "some password"
	)

	driver, err := createDriverWithAddressResolver("neo4j://x.acme.com", username, password,
		neo4j.NewServerAddress("a.acme.com", "7676"),
		neo4j.NewServerAddress("b.acme.com", "8787"),
		neo4j.NewServerAddress("c.acme.com", "9898"))
	if err != nil {
		return err
	}
	defer driver.Close()

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.Run("CREATE (n:Person { name: $name})", map[string]interface{}{"name": name})
	if err != nil {
		return err
	}

	_, err = result.Consume()
	if err != nil {
		return err
	}

	return nil
}

// end::config-custom-resolver[]

// tag::config-connection-pool[]
func createDriverWithCustomizedConnectionPool(uri, username, password string) (neo4j.Driver, error) {
	return neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""), func(config *neo4j.Config) {
		config.MaxConnectionLifetime = 30 * time.Minute
		config.MaxConnectionPoolSize = 50
		config.ConnectionAcquisitionTimeout = 2 * time.Minute
	})
}

// end::config-connection-pool[]

// tag::config-connection-timeout[]
func createDriverWithConnectionTimeout(uri, username, password string) (neo4j.Driver, error) {
	return neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""), func(config *neo4j.Config) {
		config.SocketConnectTimeout = 15 * time.Second
	})
}

// end::config-connection-timeout[]

// tag::config-max-retry-time[]
// This driver is used to run queries, needs actual TLS configuration as well.
func createDriverWithMaxRetryTime(uri, username, password string) (neo4j.Driver, error) {
	return neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""), func(config *neo4j.Config) {
		config.MaxTransactionRetryTime = 15 * time.Second
	})
}

// end::config-max-retry-time[]

// tag::service-unavailable[]
func createItem(driver neo4j.Driver) error {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("CREATE (a:Item)", nil)
		if err != nil {
			return nil, err
		}

		return result.Consume()
	})

	return err
}

// end::service-unavailable[]

func countNodes(driver neo4j.Driver, label string, property string, value string) (int64, error) {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	result, err := session.Run(fmt.Sprintf("MATCH (a:%s {%s: $value}) RETURN count(a)", label, property), map[string]interface{}{"value": value})
	if err != nil {
		return -1, err
	}

	if result.Next() {
		return result.Record().Values[0].(int64), nil
	}

	return -1, errors.New("expected at least one record")
}

// tag::session[]
func addPersonInSession(driver neo4j.Driver, name string) error {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.Run("CREATE (a:Person {name: $name})", map[string]interface{}{"name": name})
	if err != nil {
		return err
	}

	if _, err = result.Consume(); err != nil {
		return err
	}

	return nil
}

// end::session[]

// tag::autocommit-transaction[]
func addPersonInAutoCommitTx(driver neo4j.Driver, name string) error {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.Run("CREATE (a:Person {name: $name})", map[string]interface{}{"name": name})
	if err != nil {
		return err
	}

	if _, err = result.Consume(); err != nil {
		return err
	}

	return nil
}

// end::autocommit-transaction[]

// tag::transaction-function[]
func addPersonInTxFunc(driver neo4j.Driver, name string) error {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("CREATE (a:Person {name: $name})", map[string]interface{}{"name": name})
		if err != nil {
			return nil, err
		}

		return result.Consume()
	})

	return err
}

// end::transaction-function[]

// tag::transaction-timeout-config[]
func configTxTimeout(driver neo4j.Driver, name string) error {
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("CREATE (a:Person {name: $name})", map[string]interface{}{"name": name})
		if err != nil {
			return nil, err
		}

		return result.Consume()
	}, neo4j.WithTxTimeout(5*time.Second))

	return err
}

// end::transaction-timeout-config[]

// tag::transaction-metadata-config[]
func configTxMetadata(driver neo4j.Driver, name string) error {
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("CREATE (a:Person {name: $name})", map[string]interface{}{"name": name})
		if err != nil {
			return nil, err
		}

		return result.Consume()
	}, neo4j.WithTxMetadata(map[string]interface{}{"applicationId": 123}))

	return err
}

// end::transaction-metadata-config[]

// tag::pass-bookmarks[]
func addCompanyTxFunc(name string) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		return tx.Run("CREATE (a:Company {name: $name})", map[string]interface{}{"name": name})
	}
}

func addPersonTxFunc(name string) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		return tx.Run("CREATE (a:Person {name: $name})", map[string]interface{}{"name": name})
	}
}

func employTxFunc(person string, company string) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		return tx.Run(
			"MATCH (person:Person {name: $personName}) "+
				"MATCH (company:Company {name: $companyName}) "+
				"CREATE (person)-[:WORKS_FOR]->(company)", map[string]interface{}{"personName": person, "companyName": company})
	}
}

func makeFriendTxFunc(person1 string, person2 string) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		return tx.Run(
			"MATCH (a:Person {name: $name1}) "+
				"MATCH (b:Person {name: $name2}) "+
				"MERGE (a)-[:KNOWS]->(b)", map[string]interface{}{"name1": person1, "name2": person2})
	}
}

func printFriendsTxFunc() neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (a)-[:KNOWS]->(b) RETURN a.name, b.name", nil)
		if err != nil {
			return nil, err
		}

		for result.Next() {
			fmt.Printf("%s knows %s\n", result.Record().Values[0], result.Record().Values[1])
		}

		return result.Consume()
	}
}

func addAndEmploy(driver neo4j.Driver, person string, company string) (string, error) {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	if _, err := session.WriteTransaction(addCompanyTxFunc(company)); err != nil {
		return "", err
	}
	if _, err := session.WriteTransaction(addPersonTxFunc(person)); err != nil {
		return "", err
	}
	if _, err := session.WriteTransaction(employTxFunc(person, company)); err != nil {
		return "", err
	}

	return session.LastBookmark(), nil
}

func makeFriend(driver neo4j.Driver, person1 string, person2 string, bookmarks ...string) (string, error) {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, Bookmarks: bookmarks})
	defer session.Close()

	if _, err := session.WriteTransaction(makeFriendTxFunc(person1, person2)); err != nil {
		return "", err
	}

	return session.LastBookmark(), nil
}

func addEmployAndMakeFriends(driver neo4j.Driver) error {
	var bookmark1, bookmark2, bookmark3 string
	var err error

	if bookmark1, err = addAndEmploy(driver, "Alice", "Wayne Enterprises"); err != nil {
		return err
	}

	if bookmark2, err = addAndEmploy(driver, "Bob", "LexCorp"); err != nil {
		return err
	}

	if bookmark3, err = makeFriend(driver, "Bob", "Alice", bookmark1, bookmark2); err != nil {
		return err
	}

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead, Bookmarks: []string{bookmark1, bookmark2, bookmark3}})
	defer session.Close()

	if _, err = session.ReadTransaction(printFriendsTxFunc()); err != nil {
		return err
	}

	return nil
}

// end::pass-bookmarks[]

// tag::read-write-transaction[]
func addPersonNodeTxFunc(name string) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("CREATE (a:Person {name: $name})", map[string]interface{}{"name": name})
		if err != nil {
			return nil, err
		}

		return result.Consume()
	}
}

func matchPersonNodeTxFunc(name string) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (a:Person {name: $name}) RETURN id(a)", map[string]interface{}{"name": name})
		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, errors.New("one record was expected")
	}
}

func addPersonNode(driver neo4j.Driver, name string) (int64, error) {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	if _, err := session.WriteTransaction(addPersonNodeTxFunc(name)); err != nil {
		return -1, err
	}

	var id interface{}
	var err error
	if id, err = session.ReadTransaction(matchPersonNodeTxFunc(name)); err != nil {
		return -1, err
	}

	return id.(int64), nil
}

// end::read-write-transaction[]

func TestExamplesDatabaseSelection(t *testing.T) {
	driver := dbserver.GetDbServer().Driver()
	defer driver.Close()
	// tag::database-selection[]
	session := driver.NewSession(neo4j.SessionConfig{DatabaseName: "example"})
	// end::database-selection[]
	defer session.Close()
}

// tag::result-consume[]
func getPeople(driver neo4j.Driver) ([]string, error) {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	people, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		var list []string

		result, err := tx.Run("MATCH (a:Person) RETURN a.name ORDER BY a.name", nil)
		if err != nil {
			return nil, err
		}

		for result.Next() {
			list = append(list, result.Record().Values[0].(string))
		}

		if err = result.Err(); err != nil {
			return nil, err
		}

		return list, nil
	})
	if err != nil {
		return nil, err
	}

	return people.([]string), nil
}

// end::result-consume[]

// tag::result-retain[]
func addPersonsAsEmployees(driver neo4j.Driver, companyName string) (int, error) {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	persons, err := neo4j.Collect(session.Run("MATCH (a:Person) RETURN a.name AS name", nil))
	if err != nil {
		return 0, err
	}

	employees := 0
	for _, person := range persons {
		_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return tx.Run("MATCH (emp:Person {name: $person_name}) "+
				"MERGE (com:Company {name: $company_name}) "+
				"MERGE (emp)-[:WORKS_FOR]->(com)", map[string]interface{}{"person_name": person.Values[0], "company_name": companyName})
		})
		if err != nil {
			return 0, err
		}

		employees++
	}

	return employees, nil
}

// end::result-retain[]

func echo(session neo4j.Session, value interface{}) (neo4j.Record, error) {
	record, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("RETURN $value as fieldName", map[string]interface{}{"value": value})

		if err != nil {
			return neo4j.Record{}, err
		}

		if result.Next() {
			return *result.Record(), nil
		}

		return neo4j.Record{}, result.Err()

	})

	return record.(neo4j.Record), err
}
