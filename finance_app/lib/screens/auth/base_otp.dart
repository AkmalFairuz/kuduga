import 'dart:io';

import 'package:finance_app/model/otp.dart';
import 'package:finance_app/service/user_service.dart';
import 'package:finance_app/styles/color_helper.dart';
import 'package:finance_app/widgets/text/text_app_bar.dart';
import 'package:flutter/material.dart';
import 'package:url_launcher/url_launcher_string.dart';

import '../../utils/alert.dart';
import '../../widgets/button/button.dart';
import '../../widgets/input/otp_input.dart';

class BaseOTP extends StatefulWidget {
  const BaseOTP(
      {Key? key,
      required this.title,
      required this.subtitle,
      required this.description,
      required this.model,
      this.onSubmit})
      : super(key: key);

  final String title;
  final String subtitle;
  final String description;
  final OTPModel model;
  final Function(String)? onSubmit;

  @override
  State<BaseOTP> createState() => _BaseOTPState();
}

class _BaseOTPState extends State<BaseOTP> {
  String _otpCode = "";
  int lastResend = 0;

  @override
  void initState() {
    super.initState();
    lastResend = DateTime.now().millisecondsSinceEpoch;
  }

  void _onResendOTP() async {
    if (DateTime.now().millisecondsSinceEpoch - lastResend < 45000) {
      int diff = 45 -
          ((DateTime.now().millisecondsSinceEpoch - lastResend) / 1000).round();
      Alert.message("Tunggu ${diff} detik untuk mengirim ulang OTP");
      return;
    }
    lastResend = DateTime.now().millisecondsSinceEpoch;
    Alert.withLoading((done) async {
      await UserService.resendOTP(widget.model);
      done();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: TextAppBar(widget.title),
        body: Column(
          children: [
            Expanded(
                child: ListView(
              children: [
                const SizedBox(height: 80),
                Container(
                    padding: const EdgeInsets.symmetric(horizontal: 24),
                    constraints: const BoxConstraints(
                      maxWidth: 300,
                    ),
                    child: Column(
                      children: [
                        Text(widget.subtitle,
                            textAlign: TextAlign.center,
                            style: const TextStyle(
                                fontWeight: FontWeight.bold, fontSize: 26)),
                        const SizedBox(height: 32),
                        Text(widget.description,
                            style: TextStyle(color: ColorHelper.bg700(context)),
                            textAlign: TextAlign.center),
                        const SizedBox(height: 4),
                        TextButton(
                            onPressed: () {
                              if (Platform.isIOS) {
                                launchUrlString("message://");
                              } else {
                                launchUrlString("googlegmail://");
                              }
                            },
                            child: const Text("Buka Email"))
                      ],
                    )),
                const SizedBox(height: 16),
                OTPInput(
                    initialFocus: true,
                    onChanged: (otp) {
                      _otpCode = otp;
                    }),
                const SizedBox(height: 64),
                Column(
                  children: [
                    Text(
                      "Belum menerima email?",
                      style: TextStyle(color: ColorHelper.bg700(context)),
                    ),
                    TextButton(
                        onPressed: _onResendOTP,
                        child: const Text("Kirim ulang OTP"))
                  ],
                )
              ],
            )),
            Container(
              padding: const EdgeInsets.all(16),
              child: Button(
                onPressed: () {
                  if (widget.onSubmit != null) {
                    widget.onSubmit!(_otpCode);
                  } else {
                    Navigator.of(context).pop(_otpCode);
                  }
                },
                width: double.infinity,
                child: const Text("Lanjutkan"),
              ),
            )
          ],
        ));
  }
}
