class TransferUserInfoModel {
  String username;
  String name;
  String email;

  TransferUserInfoModel(
      {required this.username, required this.name, required this.email});

  factory TransferUserInfoModel.fromJson(Map<String, dynamic> json) =>
      TransferUserInfoModel(
          username: json['username'], name: json['name'], email: json['email']);
}

class TransferModel {
  int id;
  bool isSend;
  TransferUserInfoModel sender;
  TransferUserInfoModel receiver;
  int amount;
  String note;
  int createdAt;

  TransferModel(
      {required this.id,
      required this.isSend,
      required this.sender,
      required this.receiver,
      required this.amount,
      required this.note,
      required this.createdAt});

  factory TransferModel.fromJson(Map<String, dynamic> json) => TransferModel(
      id: json["id"],
      isSend: json["isSend"],
      sender: TransferUserInfoModel.fromJson(json["sender"]),
      receiver: TransferUserInfoModel.fromJson(json["receiver"]),
      amount: json["amount"],
      note: json["note"],
      createdAt: json["createdAt"]);
}
