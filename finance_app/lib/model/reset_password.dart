class ResetPasswordModel {
  int resetPasswordId;
  String token;

  ResetPasswordModel({required this.resetPasswordId, required this.token});

  factory ResetPasswordModel.fromJson(Map<String, dynamic> json) =>
      ResetPasswordModel(
          resetPasswordId: json["resetPasswordId"], token: json["token"]);
}
