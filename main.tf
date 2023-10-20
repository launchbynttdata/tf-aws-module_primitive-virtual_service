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

resource "aws_appmesh_virtual_service" "this" {
  name      = var.name
  mesh_name = var.app_mesh_name

  spec {
    provider {
      dynamic "virtual_node" {
        for_each = length(var.virtual_node_name) > 0 ? [1] : []
        content {
          virtual_node_name = var.virtual_node_name
        }

      }

      dynamic "virtual_router" {
        for_each = length(var.virtual_router_name) > 0 ? [1] : []
        content {
          virtual_router_name = var.virtual_router_name
        }
      }
    }
  }

  tags = local.tags
}
