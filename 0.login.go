package main

import (
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Write(SignInPage(SessionDeviceMode(r)))
		return
	}

	user := r.FormValue(`email`)
	pass := r.FormValue(`password`)
	info := FindUserInfo(user, pass)
	if info.login == 0 {
		http.Error(w, `Invalid email or password`, http.StatusUnauthorized)
		return
	}

	state := GetState(r)
	state.user = info
	SetState(r, state)

	http.Redirect(w, r, `/`, http.StatusSeeOther)
}

func SignOutHandler(w http.ResponseWriter, r *http.Request) {
	DestroySession(r)
	ClearSessionCookie(w)
	http.Redirect(w, r, `/signin`, http.StatusSeeOther)
}

func SignInPage(mode string) []byte {
	script := DeviceConfirmHeadScript(mode)
	page := strings.Replace(string(getLogin), `</head>`, "\t<script>\n"+script+"\n\t</script>\n</head>", 1)
	return []byte(page)
}

func FindUserInfo(user, pass string) UserInfo_t {
	var greet, hash string
	row := App.DB.CallRow(`klec_account_hash_get`, user).Scan(&greet, &hash)
	if row.HasError() || len(hash) == 0 { return UserInfo_t{} }
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)) != nil { return UserInfo_t{} }

	info := UserInfo_t{ 1, greet, user }
	out, ok := FindUserInfoByGreet(greet)
	if !ok { return info }

	if out.login == 0 { out.login = 1 }
	if out.greet == `` { out.greet = greet }
	if out.email == `` { out.email = user }
	return out
}

func FindUserInfoByGreet(greet string) (UserInfo_t, bool) {
	rows := App.DB.Call(`klec_account_get`, 0)
	if rows.HasError() { return UserInfo_t{}, false }
	defer rows.Close()

	for rows.Next() {
		var x UserInfo_t
		var crypt, created, updated string
		var editor, active bool
		if e := rows.Scan(&x.login, &x.greet, &x.email, &crypt, &editor, &active, &created, &updated); e != nil {
			return UserInfo_t{}, false
		}
		if x.greet != greet { continue }
		return x, true
	}

	return UserInfo_t{}, false
}

var getLogin = []byte(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="icon" href="/static/favicon.ico" type="image/x-icon">
    <link rel="stylesheet" type="text/css" href="/static/css/login.css">
    <title>Sign in - ` + SiteName + `</title>
</head>
<body>
    <main class="login-page">
        <section class="login-panel">
            <h1 class="login-title">Sign in</h1>
            <p class="login-subtitle">Price Machine 04. May 2026</p>
            <form action="/signin" method="post">
                <div class="login-field">
                    <label>Email or Greeting Name</label>
                    <input type="text" name="email" required>
                </div>

                <div class="login-field">
                    <label>Password</label>
                    <input type="password" name="password" required>
                </div>

                <button type="submit">Sign in</button>
            </form>
        </section>
    </main>
</body>
</html>`)
