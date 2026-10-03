class DepositDetailedModel {
  int id;
  String? layoutName;
  String paymentMethod;
  int amount;
  int fee;
  int status;
  String? paymentUrl;
  Map<String, dynamic> paymentData;
  String description;
  int createdAt;
  int expiredAt;
  List<DepositTrack> tracks;

  DepositDetailedModel.fromJson(Map<String, dynamic> json)
      : id = json['id'],
        layoutName = json['layoutName'],
        paymentMethod = json['paymentMethod'],
        amount = json['amount'],
        fee = json['fee'],
        status = json['status'],
        paymentData = json['paymentData'],
        description = json['description'],
        createdAt = json['createdAt'],
        expiredAt = json['expiredAt'],
        paymentUrl = json['paymentUrl'],
        tracks = (json['tracks'] as List<dynamic>)
            .map((e) => DepositTrack.fromJson(e))
            .toList();

  // This method only works for bank transfer payment method
  String? getBankAccountNumber() {
    return paymentData['accountNumber'];
  }

  // This method only works for bank transfer payment method
  String? getBankAccountName() {
    return paymentData['accountName'];
  }

  bool isSuccess() {
    return status == 1;
  }
}

class DepositModel {
  int id;
  String paymentMethod;
  int amount;
  int fee;
  int status;
  int createdAt;

  DepositModel.fromJson(Map<String, dynamic> json)
      : id = json['id'],
        paymentMethod = json['paymentMethod'],
        amount = json['amount'],
        fee = json['fee'],
        status = json['status'],
        createdAt = json['createdAt'];
}

class DepositTrack {
  int? newStatus;
  String description;
  int createdAt;

  DepositTrack.fromJson(Map<String, dynamic> json)
      : newStatus = json['newStatus'],
        description = json['description'],
        createdAt = json['createdAt'];
}

class DepositStatus {
  static const int waiting = 0;
  static const int success = 1;
  static const int failed = 2;

  static String displayStatus(int status) {
    switch (status) {
      case waiting:
        return "Menunggu Pembayaran";
      case success:
        return "Sukses";
      case failed:
        return "Gagal";
      default:
        return "Unknown";
    }
  }
}
