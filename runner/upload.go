package main

import (
	"fmt"
	"net/http"
	"os"

	u "github.com/mishazawa/steve/utils"
)

func Upload(root, endpoint string) error {
	listDir, err := u.GetPaths(root)
	if err != nil {
		return err
	}

	client := &http.Client{}

	u.ProcessList(listDir, func(errChan chan error, fp string) {
		file, err := os.Open(fp)
		if err != nil {
			errChan <- err
			return
		}
		defer file.Close()

		req, err := http.NewRequest("POST", endpoint, file)
		req.Header.Add("X-Filepath", fp)
		fmt.Println(fp)
		res, err := client.Do(req)

		if err != nil {
			errChan <- err
			return
		}

		if res.StatusCode != http.StatusOK {
			errChan <- fmt.Errorf("%s %s", fp, res.Status)
			return
		}

		errChan <- nil
	})

	return nil
}
