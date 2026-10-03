class OTPModel {
  int id;
  String token;
  String? code;

  OTPModel({
    required this.id,
    required this.token,
    this.code,
  });

  factory OTPModel.fromJson(Map<String, dynamic> json) {
    return OTPModel(
      id: json['otpId'],
      token: json['token'],
      code: json['code'],
    );
  }

  static bool isValidCode(String code) {
    return code.length == 6 && int.tryParse(code) != null;
  }
}
