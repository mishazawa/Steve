package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	u "github.com/mishazawa/steve/utils"
)

func Download(destination, endpoint string) error {
	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/ls", endpoint), nil)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var files []string

	err = json.Unmarshal(body, &files)
	if err != nil {
		return err
	}

	return u.ProcessList(files, func(errChan chan error, fp string) {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/static/%s", endpoint, fp), nil)
		if err != nil {
			errChan <- err
			return
		}

		res, err := client.Do(req)
		if err != nil {
			errChan <- err
			return
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			errChan <- fmt.Errorf("%s %s", fmt.Sprintf("%s/static/%s", endpoint, fp), res.Status)
			return
		}

		file, err := os.Create(fp)
		if err != nil {
			errChan <- err
			return
		}

		_, err = io.Copy(file, res.Body)
		if err != nil {
			errChan <- err
			return
		}

		errChan <- nil
	})
}
