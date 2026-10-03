import 'package:finance_app/utils/meta.dart';

import 'otp.dart';

class RegisterModel {
  String displayName;
  String username;
  String email;
  String password;
  OTPModel? otp;

  RegisterModel(
      {required this.displayName,
      required this.username,
      required this.email,
      required this.password,
      this.otp});

  void validate() {
    if (displayName.isEmpty) {
      throw RegisterException("Nama tidak boleh kosong");
    }
    if (username.isEmpty) {
      throw RegisterException("Username tidak boleh kosong");
    }
    if (email.isEmpty) {
      throw RegisterException("Email tidak boleh kosong");
    }
    if (!Meta.isValidEmail(email)) {
      throw RegisterException("Email tidak valid");
    }
    if (password.isEmpty) {
      throw RegisterException("Password tidak boleh kosong");
    }
  }
}

class RegisterException implements Exception {
  final String message;

  RegisterException(this.message);

  @override
  String toString() {
    return message;
  }
}
