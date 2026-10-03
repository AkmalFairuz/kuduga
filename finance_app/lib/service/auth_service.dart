import 'package:finance_app/model/otp.dart';
import 'package:finance_app/model/register.dart';
import 'package:finance_app/server/server.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:google_sign_in/google_sign_in.dart';

import '../model/auth_details.dart';
import '../utils/json.dart';

class AuthService {
  static Future<String> login(String usernameOrEmail, String password) {
    return Server.post('/auth/login',
        body: {'username': usernameOrEmail, 'password': password}).then(
      (resp) {
        return jsonDec(resp)['token'];
      },
    );
  }

  static Future<OTPModel> getRegisterOTP(RegisterModel model) async {
    final body = {
      'name': model.username,
      'displayName': model.displayName,
      'email': model.email,
      'password': model.password,
    };
    if (Meta.udid != null) {
      body['deviceUniqueId'] = Meta.udid!;
    }
    final resp = await Server.post('/user/registerOTP', body: body);
    return OTPModel.fromJson(jsonDec(resp));
  }

  static Future<String> register(RegisterModel model) async {
    final body = {
      'name': model.username,
      'displayName': model.displayName,
      'email': model.email,
      'password': model.password,
      'otp.id': model.otp!.id.toString(),
      'otp.token': model.otp!.token,
      'otp.code': model.otp!.code,
    };
    if (Meta.udid != null) {
      body['deviceUniqueId'] = Meta.udid!;
    }
    var resp = await Server.post('/user/register', body: body);
    return jsonDec(resp)["token"];
  }

  static Future<AuthDetailsModel> getAuthDetails() {
    return Server.get('/auth/details').then(
      (resp) {
        return AuthDetailsModel.fromJson(jsonDec(resp));
      },
    );
  }

  static Future<void> extendToken() async {
    await Server.post("/auth/extendToken");
  }

  static Future<void> logout() async {
    final gauth = GoogleSignIn();
    if (await gauth.isSignedIn()) {
      await gauth.signOut();
    }
    await Server.post('/auth/logout');
  }

  static Future<String> signUpWithGoogle() async {
    GoogleSignIn oauth = GoogleSignIn(
      scopes: ['email', 'profile'],
    );
    try {
      final result = await oauth.signIn();
      if (result == null) {
        throw Exception("Daftar dengan google dibatalkan");
      }
      final googleKey = await result.authentication;
      final body = {
        'idToken': googleKey.idToken,
      };
      if (Meta.udid != null) {
        body['deviceUniqueId'] = Meta.udid!;
      }
      final resp = await Server.post('/user/registerWithGoogle', body: body);
      return jsonDec(resp)["token"];
    } finally {
      oauth.signOut();
    }
  }

  static Future<String> signInWithGoogle() async {
    GoogleSignIn oauth = GoogleSignIn(
      scopes: ['email', 'profile'],
    );
    try {
      final result = await oauth.signIn();
      if (result == null) {
        throw Exception("Masuk dengan google dibatalkan");
      }
      final googleKey = await result.authentication;
      final resp = await Server.post('/auth/loginWithGoogle', body: {
        'idToken': googleKey.idToken,
      });
      return jsonDec(resp)["token"];
    } finally {
      oauth.signOut();
    }
  }
}
