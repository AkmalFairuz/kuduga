import 'dart:typed_data';

import 'package:finance_app/model/kyc.dart';
import 'package:finance_app/model/reset_password.dart';
import 'package:http/http.dart';

import '../model/otp.dart';
import '../model/referral.dart';
import '../server/server.dart';
import '../utils/json.dart';
import '../utils/meta.dart';

class UserService {
  static Future<void> updateDisplayName(String displayName) {
    return Server.post('/user/updateDisplayName',
        body: {'displayName': displayName});
  }

  static Future<OTPModel> requestUpdateEmailOTP(String email, String password) {
    return Server.post('/user/updateEmailOTP',
        body: {'email': email, 'password': password}).then(
      (resp) {
        return OTPModel.fromJson(jsonDec(resp));
      },
    );
  }

  static Future<String> updateEmail(String email, OTPModel otp) async {
    var resp = await Server.post('/user/updateEmail', body: {
      'email': email,
      'otp.id': otp.id.toString(),
      'otp.token': otp.token,
      'otp.code': otp.code,
    });
    return jsonDec(resp)["message"];
  }

  static Future<void> updatePassword(
      String currentPassword, String newPassword) async {
    await Server.post('/user/updatePassword', body: {
      'currentPassword': currentPassword,
      'newPassword': newPassword,
    });
  }

  static Future<String> createPin(String pin) async {
    var resp = await Server.post("/user/pin/create", body: {"pin": pin});
    return jsonDec(resp)["message"];
  }

  static Future<String> updatePin(String pin, String newPin) async {
    var resp = await Server.post("/user/pin/update",
        body: {"pin": pin, "newPin": newPin});
    return jsonDec(resp)["message"];
  }

  static Future<String> deletePin(String pin) async {
    var resp = await Server.post("/user/pin/delete", body: {"pin": pin});
    return jsonDec(resp)["message"];
  }

  static Future<bool> validatePin(String pin) async {
    var resp = await Server.post("/user/pin/validate", body: {"pin": pin});
    return jsonDec(resp)["valid"];
  }

  static Future<OTPModel> recoveryPinOTP() async {
    var resp = await Server.post("/user/pin/recoveryOTP");
    return OTPModel.fromJson(jsonDec(resp));
  }

  static Future<String> recoveryPin(OTPModel model) async {
    var resp = await Server.post("/user/pin/recovery", body: {
      "id": model.id.toString(),
      "token": model.token,
      "code": model.code
    });
    return jsonDec(resp)["message"];
  }

  static Future<OTPModel> resetPasswordOTP(String email) async {
    var resp =
        await Server.post("/user/resetPasswordOTP", body: {"email": email});
    return OTPModel.fromJson(jsonDec(resp)["otp"]);
  }

  static Future<ResetPasswordModel> resetPasswordVerifyOTP(
      String email, OTPModel otp) async {
    var resp = await Server.post("/user/resetPasswordVerifyOTP", body: {
      "email": email,
      "otp.id": otp.id.toString(),
      "otp.token": otp.token,
      "otp.code": otp.code
    });
    return ResetPasswordModel.fromJson(jsonDec(resp));
  }

  static Future<String> resetPassword(
      ResetPasswordModel model, String newPassword) async {
    var resp = await Server.post("/user/resetPassword", body: {
      "resetPasswordId": model.resetPasswordId.toString(),
      "token": model.token,
      "newPassword": newPassword
    });
    return jsonDec(resp)["message"];
  }

  static Future<KycStatusModel> getKycStatus() async {
    var resp = await Server.get("/user/kyc/status");
    return KycStatusModel.fromJson(jsonDec(resp));
  }

  static Future<String> submitKyc(
      {required String documentId,
      required String fullName,
      required Uint8List ktpData}) async {
    var resp = await Server.post("/user/kyc/request", body: {
      "documentId": documentId,
      "fullName": fullName
    }, files: [
      MultipartFile.fromBytes("document", ktpData, filename: "ktp.png")
    ]);
    return jsonDec(resp)["message"];
  }

  static Future<String> resendOTP(OTPModel model) async {
    var resp = await Server.post("/user/resendOTP",
        body: {"id": model.id.toString(), "token": model.token});
    return jsonDec(resp)["message"];
  }

  static Future<String> useReferral(String referralCode) async {
    var resp = await Server.post("/user/referral/use",
        body: {"referralCode": referralCode});
    return jsonDec(resp)["message"];
  }

  static Future<ReferralInfoModel> getReferralInfo() async {
    var resp = await Server.get("/user/referral/info");
    return ReferralInfoModel.fromJson(jsonDec(resp));
  }

  static Future<void> updateDeviceUniqueID() async {
    final deviceUniqueID = Meta.udid;
    if (deviceUniqueID == null) {
      return;
    }
    await Server.post('/user/updateDeviceUniqueID',
        body: {'deviceUniqueID': deviceUniqueID});
  }
}
