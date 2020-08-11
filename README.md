# Go-Blog

A simple Go-Blog project realized with [Revel](https://github.com/revel/revel).
It uses [bbolt](https://github.com/etcd-io/bbolt) as persistence and [storm](https://github.com/asdine/storm) as ORM.

This example aims at providing a simple application core for beginners with Go and Revel. You could think of it as a brief tutorial.
There are a lot of things you could do better, smarter and more elegant how ever the application tries to give a brief overview with a small codebase about certain ways on how to implement a login, persistence, input-validation, serversiderendering and a simple API with the above mentioned frameworks / libraries.

It is based on the Revel [booking](https://github.com/revel/examples/tree/master/booking) example and tries to simplify that a little bit.
**The explanation below might contain errors or wrong information. I'm still a beginner with Revel and Go.** 

Note that the HTML and CSS is kept very basic since this project aims not at getting used productive in anyway.

## How to start

Clone the repository and navigate into the folder above the project.
You can then do `revel run .\rala-blog\` 

The application is then available at `localhost:9000`

## Architecture

The app consists of three main components:

- The Authentication Controller
- The "App" which is responsible for creating entries for our blog
- The "Home" view which is the default view and lists our blog entries

Additionally the app has a database service which is responsible for saving and reading from out `bbolt database` with `storm`.

We have the following routes in our routes file available:

```
GET     /                                       Home.Index
GET     /PostEntry                              App.PostEntry
POST    /PostEntry                              App.ReceiveEntry

GET     /Login                                  Authentication.Login
GET     /Logout                                 Authentication.Logout
POST    /Login                                  Authentication.ReceiveLogin
```

### Home Controller

Lets start with our index route `Home.Index`:

```go
type Home struct {
	*revel.Controller
}

func (c Home) Index() revel.Result {
	entries := services.GetAllEntries()
	return c.Render(entries)
}
```

There is not much happening here. We are fetching all entries from our database service and inject it into our template rendering fuction.
On the HTML template part the following piece of code is then responsible for rendering the entries:

```html
<div class="row">
  <div class="leftcolumn">
    {{range $i, $entry := .entries}}
    <div class="card">
      <h2> {{$entry.Name}} </h2>
      <h5> {{$entry.CreatedAt.Format "2006 Jan 02"}} </h5>
      <hr>
      <p> {{$entry.Description}} </p>
    </div>
    {{end}}
  </div>
</div>
```

We are iterating through all entries and create a `div` for each. 
Additionally we are formatting the date of our entries. You can read about formatting dates in Go in the official documentation.

### App Controller

Let's have a look at the functions inside the PostEntry view (GET & POST).
Similiar to the Home-View we have a simple function which is responsible for rendering the HTML template:

```go
func (c App) PostEntry() revel.Result {
	return c.Render()
}
```

Note that our function now retrieves `c App` as parameter and not `c Home` anymore.
You have to make sure that your templates are under `views/App/`. Otherwise Revel will not find the template to render.

Taking a look at the template:

```html
...
<form method="POST" id="formLogin" action="/PostEntry">

  <input type="text" name="name" placeholder="Name">
  <br>
  <input type="textarea" name="description" placeholder="Content">
  <br>

  <p class="buttons">
    <input type="submit" value="Confirm" name="confirm">
  </p>
</form>
...
```

It is a basic HTML form which does a POST to `/PostEntry`.
We have a `PostEntry` in our routes which looks like the following: 

`POST    /PostEntry                              App.ReceiveEntry`

So when submitting our form our App-Controller function ReceiveEntry will receive the POST request.
Note: This is probably not the best way to do it. I think it would be better to not include the route and define the receiver directly inside the HTML template.
For clarification reasons though I defined the POST route in my routes and do it this way. It is not very maintainable and handy though.

```go
func (c App) ReceiveEntry(name string, description string) revel.Result {

	c.Validation.Required(name).Message("An entry name is required!")
	c.Validation.Required(description).Message("An entry description is required!")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(App.PostEntry)
	}

	// Creates a model with the received fields and the current time
	var entry models.BlogEntry
	entry.Name = name
	entry.Description = description
	entry.CreatedAt = time.Now()
	services.SaveToDatabase(entry)

	c.Flash.Success("Entry created!")
	return c.Redirect(App.PostEntry)
}
```

This is the method which gets fired by submitting our form. 
It receives the named fields of our form as parameters and then checks whether they are valid.
If they are invalid it uses Revels built in flash functionality to display an error and to redirect the user again to the `PostEntry` route to fill in a valid input.
In case the input is valid it creates an entry based on the model in `app/models`. Additionally it sets the created time to now (thats the field we are formatting in HTML template before).
Afterwards it saves the created entry with our databaseservice into our bbolt-database and redirects the user to our `PostEntry` route where he can create another entry.

The last method in our App controller is the following and is related to the authentication:

```go
func (c App) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Authentication.Login())
	}
	return nil
}
```

It basically checks before each call in our App controller whether the user is logged in and in case the user is not it redirects him to the login page.
This method is one of two methods called by our interceptors defined in `controllers/init.go`:

```go
func init() {
	revel.InterceptMethod(Authentication.addUser, revel.BEFORE)
	revel.InterceptMethod(App.checkUser, revel.BEFORE)
}
```

What is important to understand is that `InterceptMethod(..)` is scoped to the controller. That means `checkUser()` is only called before requests to our App controller.
That means it is not called when accessing our `Home.Index` which leads to the fact that you can see the blog even when you are not logged in. How ever you can not access `App.PostEntry` without be logged in otherwise `checkUser()` will redirect you to the login page. If you want to read more about interceptors you can do it [here](https://revel.github.io/manual/interceptors.html).

`Authentication.addUser` is therefore only called when interacting with the Authentication controller.

### Authentication Controller

The authentication controller is the most complex part of our blog. The login is session based.
Our login form gets rendered via the following function:

```go
func (c Authentication) Login() revel.Result {
	if user := c.connected(); user != nil {
		return c.Redirect(routes.App.PostEntry())
	}
	return c.Render()
}
```

In case the user is logged in the app will redirect the user to our `PostEntry` section.
Otherwise it renders the loginform as usual.

When entering username and password in our form (similiar to our `PostEntry` form) the following method gets called:

```go
func (c Authentication) ReceiveLogin(username, password string) revel.Result {
	user := c.getUser(username)
	if user != nil {
		err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
		if err == nil {
			c.Session["user"] = username
			c.Session.SetNoExpiration()
			c.Flash.Success("Welcome, " + username)
			return c.Redirect(routes.App.PostEntry())
		}
	}

	c.Flash.Out["username"] = username
	c.Flash.Error("Login failed")
	return c.Redirect(routes.Authentication.Login())
}
```

The first thing which happens is that a function named `getUser(username)` gets called:

```go
func (c Authentication) getUser(username string) (user *models.User) {
	user = &models.User{}
	_, err := c.Session.GetInto("fulluser", user, false)

	if user.Name == username {
		return user
	}

	*user, err = services.GetSingleUserByName(username)
	if err != nil {
		fmt.Println("Failed to find user", "user", username, "error", err)
		return nil
	}
	c.Session["fulluser"] = user
	return
}
```

As the name states this function is responsible for getting the user object. Also it creates the session object.
Our database service `GetSingleUserByName(username)` fetches the user from our database layer.
In case it can't find the user it will fire an error. If it does it stores the user into our session and returns.

Now that we have a session and a user object the `ReceiveLogin` function calls `bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))`.
Since we only store hashed password we need to compare a hashed password with our entered password. Bcrypt does that for us. 
If the method does not return an error the entered password parameter and the hashed password stored in our database are equal.
Therefore the login was successful.

Next thing we do is:

```go
c.Session["user"] = username
c.Session.SetNoExpiration()
c.Flash.Success("Welcome, " + username)
return c.Redirect(routes.App.PostEntry())
```

For testing purposes we set our session expiration to none and store under the `username` under the key `user` into our session.
Afterwards we flash a successful login message and redirect to the `PostEntry` section.

The three methods we did not cover yet in our Authentication controller are:

```go
func (c Authentication) connected() *models.User {
	if c.ViewArgs["user"] != nil {
		return c.ViewArgs["user"].(*models.User)
	}
	if username, ok := c.Session["user"]; ok {
		return c.getUser(username.(string))
	}

	return nil
}
```
This function basically checks whether we are logged in or not. 
If we are not it returns `nil` otherwise it returns our `*user` object.

```go
func (c Authentication) addUser() revel.Result {
	if user := c.connected(); user != nil {
		c.ViewArgs["user"] = user
	}
	return nil
}
```
addUser is the function we call in our interceptor. As mentioned it gets called before every request to our Authentication controller.
It calls `connected()` and checks whether our session is valid. So we do not have to log in over and over again. 

The last function we are missing is the `Logout()` function.

```go
func (c Authentication) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(routes.Home.Index())
}
```

It gets called by a line in our `header.html` which all templates use:

` <a href="{{url "Authentication.Logout"}}">Logout</a> `

The session will get deleted and you get redirected to our `Home.Index()` and are logged out.

### Databaseservice

The `databaseservice` is super simple and is basically just a wrapper for the functions provided by `storm`.
We have:

- SaveToDatabase(entry models.BlogEntry)
   - Saves a model to the database
   - (Used when creating a blog entry)
- CreateBaseUser()
   - Gets called in revel.init() so that we a user to login (`init.go`)
   - Username: rala Password: 123
- GetSingleUserByName(name string)
   - Fetches a single user out of our database by name
   - (Used when logging in and comparing passwords)
- GetAllEntries()
   - Fetches all blog entries of our database
   - (Used when rendering the blog entries and injecting them into our template)

Important is that you open and close the database connection with `storm` conscientious.
In this example I open it before and close it after each access.

## Conclusion

I can recommend to take a look at the official Revel [hotel booking](https://github.com/revel/examples/tree/master/booking) example. It definitely is cleaner implementation wise as this example and covers more functions.

If you are not familiar with Go and Revel I recommened checking out the [official manual](https://revel.github.io/manual/index.html).

Also if you have improvements in whichever ways feel free to open a pull request or to open an issue.
