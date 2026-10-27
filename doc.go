// Package librusApi is a client for the Librus Synergia API (https://synergia.librus.pl).
//
// Create a client with New, log in with Client.Login and then fetch resources
// with functions such as GetGrades, GetUsers or GetUserInfo. Sessions expire
// after ~10-15 minutes; wrap calls in WithReauth to log in again automatically.
package librusApi
