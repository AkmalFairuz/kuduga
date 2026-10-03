import 'package:finance_app/model/referral.dart';
import 'package:finance_app/service/user_service.dart';
import 'package:finance_app/widgets/layout/line.dart';
import 'package:finance_app/widgets/skeleton/skeleton.dart';
import 'package:finance_app/widgets/text/text_app_bar.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:share_plus/share_plus.dart';

class ReferralInfoScreen extends StatefulWidget {
  const ReferralInfoScreen({Key? key}) : super(key: key);

  @override
  State<ReferralInfoScreen> createState() => _ReferralInfoScreenState();
}

class _ReferralInfoScreenState extends State<ReferralInfoScreen> {
  ReferralInfoModel? model;

  @override
  void initState() {
    UserService.getReferralInfo().then((v) {
      setState(() {
        model = v;
      });
    });
    super.initState();
  }

  void onShare() {
    Share.share(model!.shareText);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: const TextAppBar("Undang Teman"),
        body: Container(
          padding: const EdgeInsets.all(16),
          height: double.infinity,
          width: double.infinity,
          child: Column(
            children: [
              if (model != null)
                Text(
                  model!.message,
                ),
              const SizedBox(height: 20),
              const Text(
                "Bagikan kode referral Anda:",
              ),
              const SizedBox(height: 20),
              model != null
                  ? Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        SelectableText(model!.yourReferralCode,
                            style: TextStyle(
                                fontWeight: FontWeight.bold, fontSize: 32.sp)),
                        const SizedBox(width: 4),
                        GestureDetector(
                          onTap: onShare,
                          child: const Icon(Icons.share),
                        ),
                      ],
                    )
                  : const Skeleton(
                      width: 200,
                      height: 40,
                    ),
              const SizedBox(height: 20),
              const Text(
                "Kode referral ini harus dimasukkan saat teman Anda mendaftar.",
              ),
              const SizedBox(height: 16),
              const Line(),
              const SizedBox(height: 16),
              Text(
                  "Jumlah pengguna yang menggunakan kode referral Anda: ${model?.totalUserUsingYourReferralCode ?? ""}"),
            ],
          ),
        ));
  }
}
