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

package stats

import (
	"strconv"
	"time"

	v1 "github.com/MottainaiCI/mottainai-server/routes/schema/v1"

	agenttasks "github.com/MottainaiCI/mottainai-server/pkg/tasks"

	context "github.com/MottainaiCI/mottainai-server/pkg/context"
	database "github.com/MottainaiCI/mottainai-server/pkg/db"
	setting "github.com/MottainaiCI/mottainai-server/pkg/settings"
	macaron "gopkg.in/macaron.v1"
)

type Stats struct {
	Running   int `json:"running"`
	Waiting   int `json:"waiting"`
	Errored   int `json:"error"`
	Failed    int `json:"failed"`
	Succeeded int `json:"succeeded"`
	Total     int `json:"total_tasks"`

	CreatedDaily   map[string]int `json:"created_daily"`
	FailedDaily    map[string]int `json:"failed_daily"`
	ErroredDaily   map[string]int `json:"errored_daily"`
	SucceededDaily map[string]int `json:"succeeded_daily"`
}

func Info(ctx *context.Context, db *database.Database) {

	ctx.Invoke(func(config *setting.Config) {

		rtasks, e := db.Driver.GetTaskByStatus(db.Config, "running")
		if e != nil {
			ctx.ServerError("Failed getting stats", e)
			return
		}
		running_tasks := len(rtasks)
		wtasks, e := db.Driver.GetTaskByStatus(db.Config, "waiting")
		if e != nil {
			ctx.ServerError("Failed getting stats", e)
			return
		}
		waiting_tasks := len(wtasks)
		etasks, e := db.Driver.GetTaskByStatus(db.Config, "error")
		if e != nil {
			ctx.ServerError("Failed getting stats", e)
			return
		}
		error_tasks := len(etasks)
		ftasks, e := db.Driver.GetTaskByStatus(db.Config, "failed")
		if e != nil {
			ctx.ServerError("Failed getting stats", e)
			return
		}
		failed_tasks := len(ftasks)
		stasks, e := db.Driver.GetTaskByStatus(db.Config, "success")
		if e != nil {
			ctx.ServerError("Failed getting stats", e)
			return
		}
		succeeded_tasks := len(stasks)

		var failed = GetStats(ftasks)
		var errored = GetStats(etasks)
		var succeded = GetStats(stasks)

		atasks := db.Driver.AllTasks(config)
		var created = GetStats(atasks)

		s := &Stats{}
		total := len(atasks)
		s.Errored = error_tasks
		s.Running = running_tasks
		s.Total = total
		s.Waiting = waiting_tasks
		s.Failed = failed_tasks
		s.Succeeded = succeeded_tasks
		s.CreatedDaily = created
		s.FailedDaily = failed
		s.ErroredDaily = errored
		s.SucceededDaily = succeded

		ctx.JSON(200, s)
	})
}

func GetStats(atasks []agenttasks.Task) map[string]int {
	var created = make(map[string]int)
	for _, t := range atasks {
		//t, _ := db.Driver.GetTask(i)
		t1, _ := time.Parse(
			"20060102150405",
			t.CreatedTime)
		day := strconv.Itoa(t1.Year()) + "-" + strconv.Itoa(int(t1.Month())) + "-" + strconv.Itoa(t1.Day())
		i, ok := created[day]
		if !ok {
			created[day] = 1
		} else {
			created[day] = i + 1
		}
	}
	return created
}

func Setup(m *macaron.Macaron) {
	//reqSignIn := context.Toggle(&context.ToggleOptions{SignInRequired: true})

	m.Invoke(func(config *setting.Config) {
		m.Group(config.GetWeb().GroupAppPath(), func() {
			v1.Schema.GetStatsRoute("info").ToMacaron(m, Info)
		})
	})
}
