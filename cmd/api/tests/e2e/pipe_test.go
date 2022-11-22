package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"testing"
)

/*
TestPipeCreationFlow tests the flow involved in pipe creation
--------------------
# Tested endpoints:
---| /v1/pipe (POST)
*/
func TestPipeCreationFlow(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	t.Run("/v1/pipe", func(t *testing.T) {
		t.Run("create pipe without cover photo", func(t *testing.T) {

			resJSON := struct {
				Message string `json:"message"`
			}{}
			reqBody := []byte(`{"name": "testpipe1"}`)

			// create an empty buffer and a writer to the empty buffer
			buff := &bytes.Buffer{}
			buffWriter := multipart.NewWriter(buff)

			// create a multipart header for just the name input since
			// we're testing pipe creation without a cover_photo
			reqHeader := textproto.MIMEHeader{}
			reqHeader.Set("Content-Type", "application/json; charset=UTF-8")
			namePart, err := buffWriter.CreatePart(reqHeader)
			if err != nil {
				t.Fatalf("an error occurred %s", err)
			}
			_, err = namePart.Write(reqBody)
			if err != nil {
				t.Fatalf("an error occurred %s", err)
			}
			// close writer
			err = buffWriter.Close()
			if err != nil {
				t.Fatalf("could not close writer: %s", err)
			}

			// proceed to create the actual request to create the pipe
			req, err := http.NewRequest(http.MethodPost, "/v1/pipe/", bytes.NewReader(buff.Bytes()))
			if err != nil {
				t.Fatalf("could not create request %s", err)
			}
			req.Header.Add("Content-Type", "multipart/related; boundary="+buffWriter.Boundary())
			req = attachAuthHeader(req)
			res := executeRequest(req)
			resBody, _ := io.ReadAll(res.Body)
			json.Unmarshal(resBody, &resJSON)
			t.Log(resJSON.Message)

			checkResponseCode(t, http.StatusCreated, res.Code)
		})

		t.Run("create pipe with cover photo", func(t *testing.T) {
			// just trying to run a beautiful test of mine
		})
	})
}

/*
TestGetPipeFlow tests the flow involved in fetching a pipe
--------------------
# Tested endpoints:
---| /v1/pipe/:id (GET)
---| /v1/pipe/all (GET)
*/
func TestGetPipeFlow(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	t.Run("testing pipe fetch flow", func(t *testing.T) {
		t.Run("/v1/pipe/:id", func(t *testing.T) {
			// Run the test for getting a single pipe
		})

		t.Run("/v1/pipe/all", func(t *testing.T) {
			// Run the test for getting all the pipes for the
			// currently logged in user
		})
	})
}

/*
TestMutatePipeFlow tests the flow involved in updating/deleting a single pipe
--------------------
# Tested endpoints:
---| /v1/pipe/:id (PUT)
---| /v1/pipe/:id (DELETE)
*/
func TestMutatePipeFlow(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	t.Run("testing pipe mutation flow", func(t *testing.T) {
		t.Run("(PUT)-/v1/pipe/:id", func(t *testing.T) {
			// Run the test for updating a pipe
		})

		t.Run("(DELETE)-/v1/pipe/:id", func(t *testing.T) {
			// Run the test for deleting a pipe
		})
	})
}
