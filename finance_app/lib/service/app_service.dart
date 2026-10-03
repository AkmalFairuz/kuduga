import 'package:finance_app/server/server.dart';
import 'package:finance_app/utils/json.dart';

import '../model/app_layout_info.dart';

class AppService {
  static Future<AppVersionInfo> getVersionInfo() async {
    var resp = await Server.get("/app/version");
    return AppVersionInfo.fromJson(jsonDec(resp));
  }

  static Future<AppLayoutInfo> getLayoutInfo() async {
    var resp = await Server.get("/app/layoutInfo");
    return AppLayoutInfo.fromJson(jsonDec(resp));
  }

  static Future<List<AppBannerInfo>> getBannerInfo() async {
    var resp = await Server.get("/app/bannerInfo");
    return (jsonDec(resp) as List<dynamic>)
        .map((e) => AppBannerInfo.fromJson(e))
        .toList();
  }

  static Future<String> getPrivacyPolicy() async {
    var resp = await Server.get("/app/privacyPolicy");
    return resp.body;
  }

  static Future<String> getTermsAndConditions() async {
    var resp = await Server.get("/app/termsAndConditions");
    return resp.body;
  }
}
