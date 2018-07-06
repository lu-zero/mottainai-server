/*

Copyright (C) 2018  Ettore Di Giacinto <mudler@gentoo.org>
Credits goes also to Gogs authors, some code portions and re-implemented design
are also coming from the Gogs project, which is using the go-macaron framework
and was really source of ispiration. Kudos to them!

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.

*/

package client

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	token "github.com/MottainaiCI/mottainai-server/pkg/token"

	client "github.com/MottainaiCI/mottainai-server/pkg/client"
	database "github.com/MottainaiCI/mottainai-server/pkg/db"
	node "github.com/MottainaiCI/mottainai-server/pkg/nodes"
	user "github.com/MottainaiCI/mottainai-server/pkg/user"

	"github.com/MottainaiCI/mottainai-server/pkg/mottainai"
	s "github.com/MottainaiCI/mottainai-server/pkg/settings"
	"github.com/MottainaiCI/mottainai-server/routes"
	"github.com/go-macaron/binding"
)

func TestUpload(t *testing.T) {
	binding.MaxMemory = int64(1024 * 1024 * 1)
	s.Configuration.GenDefault()
	s.Configuration.Unmarshal()
	s.Configuration.DBPath = "./DB"

	defer os.RemoveAll(s.Configuration.DBPath)
	defer os.RemoveAll(s.Configuration.ArtefactPath)

	db := database.NewDatabase("")

	u := &user.User{}
	u.Name = "test"
	u.Password = "foo"
	u.Email = "foo@bar"
	id, err := db.InsertAndSaltUser(u)
	if err != nil {
		t.Error(err)
	}

	tok, err := token.GenerateUserToken(id)
	if err != nil {
		t.Error(err)
	}

	node := &node.Node{Key: "test"}
	nodeid, err := db.InsertNode(node)
	if err != nil {
		t.Error(err)
	}
	_, err = db.InsertToken(tok)
	if err != nil {
		t.Error(err)
	}
	server := mottainai.Classic()
	routes.SetupWebUI(server)
	go server.Start()
	time.Sleep(time.Duration(5 * time.Second))
	c := client.NewTokenClient(s.Configuration.AppURL, tok.Key)

	dat := make(map[string]interface{})

	var flagsName []string = []string{
		"script", "storage", "source", "directory", "task", "image",
		"namespace", "storage_path", "artefact_path", "tag_namespace",
		"prune", "queue", "cache_image",
	}

	for _, n := range flagsName {
		dat[n] = "test"
	}

	res, err := c.GenericForm("/api/tasks", dat)
	if err != nil {
		t.Errorf(err.Error())
	}
	tid := string(res)

	fmt.Println("Created Task:", tid)

	if tid == "0" {
		t.Fatal("Document not created")
	}

	fetcher := client.NewTokenClient(s.Configuration.AppURL, tok.Key)
	fetcher.Doc(tid)
	testFile := "integration_test.go"

	s.Configuration.AgentKey = node.Key
	fetcher.RegisterNode("foo", "bar")

	nd, err := db.GetNode(nodeid)
	if err != nil {
		t.Errorf(err.Error())
	}

	if nd.NodeID != "foo" {
		t.Error("Failed registering node", nd)
	}

	err = fetcher.UploadArtefactRetry(testFile, "/", 5)
	if err != nil {
		t.Errorf(err.Error())
	}

	file, err := os.Open(path.Join(s.Configuration.ArtefactPath, tid, testFile))
	if err != nil {
		t.Errorf(err.Error())
	}

	fi, err := file.Stat()
	if err != nil {
		t.Errorf(err.Error())
	}

	t.Log(fi.Size())
}