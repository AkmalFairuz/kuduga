import 'package:flutter_contacts/flutter_contacts.dart' as FlutterContacts;

class Contact {
  final String name;
  final String phone;
  final String sanitizedPhone;

  Contact(this.name, this.phone, this.sanitizedPhone);

  static String sanitizePhone(String phone) {
    String ret = phone.replaceAll(" ", "");
    ret = ret.replaceAll("-", "");
    if (ret.startsWith("+62")) {
      ret = ret.replaceFirst("+62", "0");
    }
    return ret;
  }

  // fromContact
  Contact.fromContact(FlutterContacts.Contact contact)
      : name = contact.displayName,
        phone = contact.phones.first.number,
        sanitizedPhone = sanitizePhone(contact.phones.first.number);

  static Future<bool> requestPermission() async {
    return await FlutterContacts.FlutterContacts.requestPermission();
  }

  static Future<List<Contact>> getContacts() async {
    List<FlutterContacts.Contact> contacts =
        await FlutterContacts.FlutterContacts.getContacts(withProperties: true);

    // filter contacts, remove contacts without phone number
    return contacts
        .where((contact) =>
            contact.phones.isNotEmpty && contact.phones.first.number.isNotEmpty)
        .map((contact) => Contact.fromContact(contact))
        .toList();
  }
}
