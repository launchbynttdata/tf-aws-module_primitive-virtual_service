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

output "virtual_service_id" {
  description = "ID of the Virtual Service."
  value       = module.appmesh_virtual_service.id
}

output "virtual_service_name" {
  description = "Name of the Virtual Service."
  value       = module.appmesh_virtual_service.name
}

output "service_arn" {
  description = "ARN of the Virtual Service"
  value       = module.appmesh_virtual_service.arn
}

output "virtual_node_name" {
  description = "Name of the Virtual Node"
  value       = module.virtual_node.name
}

output "random_int" {
  description = "Random Int postfix"
  value       = random_integer.priority.result
}

output "vpc_id" {
  description = "VPC Id for the example"
  value       = module.vpc.vpc_id
}

output "mesh_name" {
  description = "Name of the Mesh"
  value       = module.app_mesh.name
}
