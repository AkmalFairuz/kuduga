import 'package:finance_app/model/otp.dart';
import 'package:finance_app/model/register.dart';
import 'package:finance_app/screens/auth/base_otp.dart';
import 'package:finance_app/server/server.dart';
import 'package:finance_app/utils/alert.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:flutter/material.dart';

import '../../service/auth_service.dart';

class RegisterOTPScreen extends StatelessWidget {
  const RegisterOTPScreen({Key? key, required this.model}) : super(key: key);

  final RegisterModel model;

  void _onSubmit(String otpCode) async {
    if (!OTPModel.isValidCode(otpCode)) {
      Alert.message("Kode OTP tidak valid");
      return;
    }

    model.otp!.code = otpCode;

    Alert.withLoading((done) async {
      final token = await AuthService.register(model);
      Server.setToken(token);
      done();
      Meta.afterRegister();
    });
  }

  @override
  Widget build(BuildContext context) {
    return BaseOTP(
        model: model.otp!,
        title: "Daftar",
        subtitle: "Verifikasi Email",
        description:
            "Masukkan 6 digit angka (OTP) yang kami telah kami kirimkan ke ${model.email}",
        onSubmit: _onSubmit);
  }
}
