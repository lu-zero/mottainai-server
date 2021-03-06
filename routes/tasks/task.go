/*

Copyright (C) 2017-2018  Ettore Di Giacinto <mudler@gentoo.org>
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

package tasks

import (
	"github.com/MottainaiCI/mottainai-server/pkg/context"

	setting "github.com/MottainaiCI/mottainai-server/pkg/settings"
	agenttasks "github.com/MottainaiCI/mottainai-server/pkg/tasks"
	"github.com/go-macaron/binding"
	macaron "gopkg.in/macaron.v1"
)

func Setup(m *macaron.Macaron) {

	m.Invoke(func(config *setting.Config) {
		reqSignIn := context.Toggle(&context.ToggleOptions{
			SignInRequired: true,
			Config:         config,
			BaseURL:        config.GetWeb().AppSubURL})

		bind := binding.BindIgnErr

		m.Group(config.GetWeb().GroupAppPath(), func() {
			m.Get("/tasks", reqSignIn, ShowAll)
			m.Get("/tasks/my", reqSignIn, ShowMyTasks)
			m.Get("/tasks/add", reqSignIn, Add)
			m.Post("/tasks", reqSignIn, bind(agenttasks.Task{}), Create)
			m.Get("/tasks/:id", reqSignIn, DisplayTask)
			m.Get("/tasks/start/:id", reqSignIn, SendStartTask)
			m.Get("/tasks/stop/:id", reqSignIn, Stop)
			m.Get("/tasks/delete/:id", reqSignIn, Delete)
			m.Get("/tasks/clone/:id", reqSignIn, Clone)
			m.Get("/tasks/status/:status", reqSignIn, ShowTaskByStatus)

			m.Get("/pipeline/:id", reqSignIn, DisplayPipeline)
			m.Get("/pipelines", reqSignIn, DisplayAllPipelines)

			//Public
			m.Get("/tasks/display/:id", reqSignIn, DisplayTask)
			m.Get("/tasks/artefacts/:id", reqSignIn, ShowArtefacts)
		})
	})
}
