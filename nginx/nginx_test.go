package nginx_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/30x/keymaster/nginx"
	"io/ioutil"
	"os"
	"path"
	"net/http"
	"time"
)

var _ = Describe("nginx", func() {

	Describe("TestConfig", func() {

		It("should return no error from good config", func() {
			tmpfile, err := writeConf(valid_conf)
			Expect(err).NotTo(HaveOccurred())
			defer os.RemoveAll(path.Dir(tmpfile.Name()))

			err = nginx.TestConfig(tmpfile.Name())
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return an error from a bad config", func() {
			tmpfile, err := writeConf(invalid_conf)
			Expect(err).NotTo(HaveOccurred())
			defer os.RemoveAll(path.Dir(tmpfile.Name()))

			err = nginx.TestConfig(tmpfile.Name())
			Expect(err).To(HaveOccurred())
		})

		It("should return an error when warnings", func() {
			tmpfile, err := writeConf(warn_conf)
			Expect(err).NotTo(HaveOccurred())
			defer os.RemoveAll(path.Dir(tmpfile.Name()))

			err = nginx.TestConfig(tmpfile.Name())
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("control", func() {

		It("should start with a good config", func() {
			tmpFile, err := writeConf(valid_conf)
			Expect(err).NotTo(HaveOccurred())
			tmpDir := path.Dir(tmpFile.Name())
			defer os.RemoveAll(tmpDir)

			err = nginx.Start(tmpDir, tmpFile.Name())
			Expect(err).NotTo(HaveOccurred())
			defer nginx.Stop(tmpDir)

			res, _ := http.Get("http://localhost:9000")
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(Equal("Hello, world\n"))
		})

		It("should reload to a upgraded config", func() {
			tmpFile, err := writeConf(valid_conf)
			Expect(err).NotTo(HaveOccurred())
			tmpDir := path.Dir(tmpFile.Name())
			defer os.RemoveAll(tmpDir)

			err = nginx.Start(tmpDir, tmpFile.Name())
			Expect(err).NotTo(HaveOccurred())
			defer nginx.Stop(tmpDir)

			res, _ := http.Get("http://localhost:9000")
			defer res.Body.Close()
			body, _ := ioutil.ReadAll(res.Body)
			Expect(string(body)).To(Equal("Hello, world\n"))

			err = ioutil.WriteFile(tmpFile.Name(), []byte(upgraded_conf), 0644)
			Expect(err).NotTo(HaveOccurred())
			err = nginx.Reload(tmpDir, tmpFile.Name())
			Expect(err).NotTo(HaveOccurred())
			time.Sleep(150 * time.Millisecond) // give it a moment to reload

			res, _ = http.Get("http://localhost:9000")
			body, err = ioutil.ReadAll(res.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(Equal("Hello, Scott\n"))
		})

		It("should fail to reload a bad config", func() {
			tmpFile, err := writeConf(valid_conf)
			Expect(err).NotTo(HaveOccurred())
			tmpDir := path.Dir(tmpFile.Name())
			defer os.RemoveAll(tmpDir)

			err = nginx.Start(tmpDir, tmpFile.Name())
			Expect(err).NotTo(HaveOccurred())
			defer nginx.Stop(tmpDir)

			res, _ := http.Get("http://localhost:9000")
			defer res.Body.Close()
			body, _ := ioutil.ReadAll(res.Body)
			Expect(string(body)).To(Equal("Hello, world\n"))

			err = ioutil.WriteFile(tmpFile.Name(), []byte(invalid_conf), 0644)
			Expect(err).NotTo(HaveOccurred())
			err = nginx.Reload(tmpDir, tmpFile.Name())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("nginx: [emerg] no \"events\" section in configuration\n"))

			res, _ = http.Get("http://localhost:9000")
			defer res.Body.Close()
			body, err = ioutil.ReadAll(res.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(Equal("Hello, world\n"))
		})
	})
})

func writeConf(content string) (*os.File, error) {
	tmpDir, err := ioutil.TempDir("", "TestConfig")
	if err != nil {
		return nil, err
	}

	tmpFile, err := ioutil.TempFile(tmpDir, "TestConfig")
	if err != nil {
		return nil, err
	}

	_, err = tmpFile.Write([]byte(content))
	if err != nil {
		return nil, err
	}
	err = tmpFile.Close()

	return tmpFile, err
}

var valid_conf = `
events {}
http {
  server {
      listen  9000;
      server_name  localhost;
      location / {
          echo Hello, world;
	  echo_flush;
      }
  }
}`

var upgraded_conf = `
events {}
http {
  server {
      listen  9000;
      server_name  localhost;
      location / {
          echo Hello, Scott;
	  echo_flush;
      }
  }
}`

var warn_conf = `
events {}
http {
  server {
      listen  9000;
      server_name  example.com;
      location / {
      }
  }
  server {
      listen  9000;
      server_name  example.com;
      location / {
      }
  }
}`

var invalid_conf = `
http {
  server {
      listen  9000;
      server_name  example.com;
      location / {
      }
  }
}`
