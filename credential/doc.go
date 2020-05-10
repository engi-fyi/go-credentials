/*
Copyright (c) 2020 engi.fyi Contributors, All Rights Reserved.

Licensed under the MIT License (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	https://engi.fyi/mit-license/

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
Package credential provides the base implementation detail for the go-credentials library.

All of the logic for handling individual credentials is contained within this package. Before a Credential can be
created, global application-level settings must be set in a factory.Factory object. Attributes aside from Username
and Password are stored in a Profile object.
*/
package credential
