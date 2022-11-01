## auth program

1. client request a business url.
2. pass_sdk filter check who is the user.
  - if user info is found, allow to visit the business page.
  - else jump to passport web site, with params: response_type=code client_id, state (it means signature and auth server will return this without any changed), redirect_uri (pass_sdk auth-path and have one param named rd, means current url)
3. user on passport site complete sign in, then it will jump back to pass_sdk's auth-path with params: code and state.
4. pass_sdk will POST a web Authorization request with grant_type=authorization_code, client_id, code and redirect_uri (same with the redirect_uri of setp 2).
5. passport server will response a json, like:
```
{
  access_token: string;
  token_type: bearer | mac | ...;
  expires_in: int;
  refresh_token: string;
  scope: string;
}
```
