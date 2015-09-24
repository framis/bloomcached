package main

import (
	"log"
	"testing"
)

var (
	c     *Client
	err   error
	dsn   = "localhost:3333"
	items = [5]string{"a",
		"122451",
		"4293f250-a85d-3f86-8197-6146faf517aa",
		"http://www.tripadvisor.com/Hotel_Review-g56241-d7374617-Reviews-Red_Door_Bed_Breakfast-Mason_Texas.html",
		"http://www.tripadvisor.com/AllLocations-g1-Places-World.html",
	}
)

func init() {
	c, err = NewClient(dsn)
	if err != nil {
		log.Fatal(err)
	}
}

func TestAdd(t *testing.T) {
	_, err := c.Add("hello")
	checkErr(t, err, "Error while adding element 'hello'")
}

func TestTest(t *testing.T) {
	for _, item := range items {
		_, err := c.Add(item)
		checkErr(t, err, "Error while adding element"+item)
		result, err := c.Test(item)
		checkErr(t, err, "Error while adding element "+item)
		if !result {
			t.Errorf(item + " should be in the filter")
		}
	}

	result, err := c.Test("Not There")
	checkErr(t, err, "")
	if result {
		t.Errorf("'Not There' should not be in the filter")
	}
}

func checkErr(t *testing.T, err error, msg string) {
	if err != nil {
		t.Errorf(msg, err.Error())
	}
}
