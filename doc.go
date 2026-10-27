// Package librus is a client for the Librus Synergia API (https://synergia.librus.pl).
//
// Create a client with New, log in with Client.Login and then fetch resources
// with its methods, such as Client.Grades, Client.Users or Client.Me:
//
//	client := librus.New()
//	if err := client.Login(ctx, login, password); err != nil {
//		// handle error
//	}
//	grades, err := client.Grades(ctx)
//
// Librus sessions expire after ~10-15 minutes. The client remembers the
// credentials passed to Login and logs in again on its own when that happens.
package librus
