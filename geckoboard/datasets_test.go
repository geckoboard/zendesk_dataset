package geckoboard

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDatasetsCreate(t *testing.T) {
	start := time.Now()

	d := DataSet{
		ID: "foobar",
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("Expected PUT request, got %q", r.Method)
		}

		if r.URL.Path != "/datasets/foobar" {
			t.Fatalf(`Expected path to be "/datasets/foobar", got %q`, r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DataSet{
			ID: "foobar",
			Fields: Fields{
				"foo": Field{},
				"bar": Field{},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}))

	c := New(Config{URL: s.URL})

	err := d.FindOrCreate(c)
	if err != nil {
		t.Fatalf("Expected not errors, got %v", err)
	}

	if len(d.Fields) != 2 {
		t.Fatalf("Expected 2 fields, got %d", len(d.Fields))
	}

	if !d.CreatedAt.After(start) {
		t.Fatalf("CreatedAt is not set correctly")
	}

	if !d.UpdatedAt.After(start) {
		t.Fatalf("UpdatedAt is not set correctly")
	}
}

func TestDatasetsCreateError(t *testing.T) {
	d := DataSet{
		ID: "foobar",
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Error{InnerError{"FoobarError"}})
	}))

	c := New(Config{URL: s.URL})

	err := d.FindOrCreate(c)
	if err == nil {
		t.Fatalf("Expected error to occur")
	}

	if err.Error() != "FoobarError" {
		t.Fatalf("Expected error message to equal FoobarError, got %q", err.Error())
	}
}

func TestDatasetsDelete(t *testing.T) {
	d := DataSet{
		ID: "foobar",
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Fatalf("Expected DELETE request, got %q", r.Method)
		}

		if r.URL.Path != "/datasets/foobar" {
			t.Fatalf(`Expected path to be "/datasets/foobar", got %q`, r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct{}{})
	}))

	c := New(Config{URL: s.URL})

	err := d.Delete(c)
	if err != nil {
		t.Fatalf("Expected not errors, got %v", err)
	}
}

func TestDatasetsDeleteError(t *testing.T) {
	d := DataSet{
		ID: "foobar",
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Error{InnerError{"FoobarError"}})
	}))

	c := New(Config{URL: s.URL})

	err := d.Delete(c)
	if err == nil {
		t.Fatalf("Expected error to occur")
	}

	if err.Error() != "FoobarError" {
		t.Fatalf("Expected error message to equal FoobarError, got %q", err.Error())
	}
}

func TestDatasetsSendAllData(t *testing.T) {
	recs := []Record{
		{
			"a": "f",
			"b": "o",
			"c": "o",
		},
		{
			"a": "b",
			"b": "a",
			"c": "r",
		},
	}

	d := DataSet{
		ID: "foobar",
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("Expected PUT request, got %q", r.Method)
		}

		if r.URL.Path != "/datasets/foobar/data" {
			t.Fatalf(`Expected path to be "/datasets/foobar/data", got %q`, r.URL.Path)
		}

		var body struct {
			Data []Record `json:"data"`
		}

		json.NewDecoder(r.Body).Decode(&body)

		if len(body.Data) != 2 {
			t.Fatalf("Expected 2 records, got %d", len(body.Data))
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct{}{})
	}))

	c := New(Config{URL: s.URL})

	err := d.SendAll(c, recs)
	if err != nil {
		t.Fatalf("Expected not errors, got %v", err)
	}
}

func TestDatasetsSendAllError(t *testing.T) {
	d := DataSet{
		ID: "foobar",
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Error{InnerError{"FoobarError"}})
	}))

	c := New(Config{URL: s.URL})

	err := d.SendAll(c, []Record{})
	if err == nil {
		t.Fatalf("Expected error to occur")
	}

	if err.Error() != "FoobarError" {
		t.Fatalf("Expected error message to equal FoobarError, got %q", err.Error())
	}
}
