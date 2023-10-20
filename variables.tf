// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

variable "name" {
  description = "Name of the virtual service"
  type        = string
}

variable "app_mesh_name" {
  description = "Id of the App Mesh where the virtual service will reside"
  type        = string
}

variable "virtual_node_name" {
  description = "Name of the virtual Node to associate with the virtual service. Conflicts with virtual_router"
  type        = string
  default     = ""
}

variable "virtual_router_name" {
  description = "Name of the virtual Router to associate with the virtual service. Conflicts with virtual_node"
  type        = string
  default     = ""
}

variable "tags" {
  description = "A map of custom tags to be attached to this resource"
  type        = map(string)
  default     = {}
}
