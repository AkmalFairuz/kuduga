class ReferralInfoModel {
  String yourReferralCode;
  String message;
  String shareText;
  int totalUserUsingYourReferralCode;

  ReferralInfoModel(
      {required this.yourReferralCode,
      required this.shareText,
      required this.message,
      required this.totalUserUsingYourReferralCode});

  factory ReferralInfoModel.fromJson(Map<String, dynamic> json) {
    return ReferralInfoModel(
        yourReferralCode: json['yourReferralCode'],
        shareText: json['shareText'],
        message: json['message'],
        totalUserUsingYourReferralCode: json['totalUserUsingYourReferralCode']);
  }
}
