import 'package:finance_app/model/otp.dart';
import 'package:finance_app/screens/auth/base_otp.dart';
import 'package:finance_app/service/user_service.dart';
import 'package:finance_app/utils/alert.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/button/button.dart';
import 'package:finance_app/widgets/input/input.dart';
import 'package:flutter/material.dart';

import '../../widgets/text/text_app_bar.dart';

class ChangeEmailScreen extends StatelessWidget {
  const ChangeEmailScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const TextAppBar("Ubah Alamat Email"),
      body: Container(
          padding: const EdgeInsets.all(16),
          child: const _ChangeEmailContent()),
    );
  }
}

class _ChangeEmailContent extends StatefulWidget {
  const _ChangeEmailContent({Key? key}) : super(key: key);

  @override
  State<_ChangeEmailContent> createState() => _ChangeEmailContentState();
}

class _ChangeEmailContentState extends State<_ChangeEmailContent> {
  FocusNode? emailFocusNode = FocusNode();
  TextEditingController emailController = TextEditingController();
  TextEditingController passwordController = TextEditingController();
  TextEditingController otpController = TextEditingController();

  OTPModel? _otpModel;

  @override
  void initState() {
    super.initState();

    Future.delayed(const Duration(milliseconds: 200), () {
      if (!mounted) return;
      if (emailFocusNode!.hasFocus) return;
      if (emailFocusNode == null) return;
      FocusScope.of(context).requestFocus(emailFocusNode);
    });
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Input(
            focusNode: emailFocusNode,
            label: 'Alamat Email Baru',
            controller: emailController),
        const SizedBox(height: 16),
        Input(
          isPassword: true,
          label: 'Password',
          controller: passwordController,
        ),
        const SizedBox(height: 16),
        Button(
            onPressed: changeEmail,
            width: double.infinity,
            child: Text(_otpModel == null ? 'Lanjutkan' : 'Ganti Email')),
      ],
    );
  }

  void _onChangeEmail(String otpCode) async {
    if (_otpModel == null) {
      return;
    }
    _otpModel!.code = otpCode;

    Alert.withLoading((done) async {
      var resp =
          await UserService.updateEmail(emailController.text, _otpModel!);
      done();
      await Alert.message(resp);
      Screens.back();
    }, message: "Mengubah email...");
  }

  void changeEmail() {
    if (emailController.text.isEmpty) {
      Alert.message("Email harus diisi");
      return;
    }
    if (passwordController.text.isEmpty) {
      Alert.message("Password harus diisi");
      return;
    }

    Alert.withLoading((done) async {
      final otp = await UserService.requestUpdateEmailOTP(
          emailController.text, passwordController.text);
      otpController.clear();
      _otpModel = otp;
      done();
      Screens.to(BaseOTP(
          model: _otpModel!,
          title: "Ubah Alamat Email",
          subtitle: "Verifikasi",
          description:
              "Masukkan 6 digit angka (OTP) yang kami telah kami kirimkan ke ${emailController.text} untuk mengganti email pada akun anda",
          onSubmit: _onChangeEmail));
    });
  }

  @override
  void dispose() {
    // dispose the focus node
    emailFocusNode!.dispose();
    emailFocusNode = null;
    super.dispose();
  }
}
