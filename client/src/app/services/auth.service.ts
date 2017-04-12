import { Injectable } from '@angular/core';
import { AuthHttp, JwtHelper, tokenNotExpired } from 'angular2-jwt';
import { Observable } from "rxjs";
import { Http } from '@angular/http';
import 'rxjs/add/operator/map';

@Injectable()
export class AuthService {
  private loginUrl:string = "/login/";
  private signupUrl:string = "/signup/";
  private logoutUrl:string = "/logout/";
  public errorMessage: string;
  public pendingRequest: boolean;
  private jwtHelper: JwtHelper = new JwtHelper();
  public username: string;
  public userID: number;

  constructor(private authHttp: AuthHttp, private http: Http) {
    if (this.loggedIn()) {
      this.username = localStorage.getItem("username")
      this.userID = +localStorage.getItem("user_id")
    }
  }

  logout() {
    if (this.loggedIn()) {
      this.pendingRequest = true;
      return this.http.post(this.logoutUrl, JSON.stringify({"logout":true})).map(res => {
            this.pendingRequest = false;
            if (res.status < 400) {
                localStorage.removeItem("token")
                localStorage.removeItem("username")
            }
        });
    }
  }

  login(username:string, password:string):Observable<JSON> {
    let requestBody = {
      "username": username,
      "password": password
    }
    this.pendingRequest = true;
    return this.http.post(this.loginUrl, JSON.stringify(requestBody))
        .map(res => {
            this.pendingRequest = false;
            if (res.status < 400) {
              let data = res.json()
              localStorage.setItem("token", data.token)
              let tokenData = this.jwtHelper.decodeToken(data.token)
              localStorage.setItem("username", tokenData.username)
              localStorage.setItem("user_id", tokenData.user_id)
              this.username = tokenData.username
              this.userID = tokenData.user_id
              return data
            }
        });
  }
  signup(username:string, email:string, password:string):Observable<JSON> {
    let requestBody = {
      "username": username,
      "email": email,
      "password": password
    }
    this.pendingRequest = true;
    return this.http.post(this.signupUrl, JSON.stringify(requestBody))
        .map(res => {
            this.pendingRequest = false;
            if (res.status < 400) {
                let data = res.json()
                localStorage.setItem("token", data.token)
                let tokenData = this.jwtHelper.decodeToken(data.token)
                localStorage.setItem("username", tokenData.username)
                localStorage.setItem("user_id", tokenData.user_id)
                this.username = tokenData.username
                this.userID = tokenData.user_id
                return data
            }
        });
  }

  loggedIn():boolean {
    return tokenNotExpired('token')
  }

}
