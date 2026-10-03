class AuthDetailsModel {
  int userId;
  String username;
  String displayName;
  String email;
  bool hasPin;

  AuthDetailsModel({
    required this.userId,
    required this.username,
    required this.displayName,
    required this.email,
    required this.hasPin,
  });

  factory AuthDetailsModel.fromJson(Map<String, dynamic> json) {
    return AuthDetailsModel(
      userId: json['userId'],
      username: json['userName'],
      displayName: json['userDisplayName'],
      email: json['userEmail'],
      hasPin: json['hasPin'],
    );
  }
}
