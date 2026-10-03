import 'package:finance_app/utils/alert.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/layout/line.dart';
import 'package:finance_app/widgets/skeleton/skeleton.dart';
import 'package:flutter/material.dart';

import '../../model/contact.dart';
import '../../widgets/text/text_app_bar.dart';

class SelectContactScreen extends StatefulWidget {
  const SelectContactScreen({Key? key}) : super(key: key);

  @override
  State<SelectContactScreen> createState() => _SelectContactState();
}

class _SelectContactState extends State<SelectContactScreen> {
  bool _granted = false;
  List<Contact> _contacts = [];

  @override
  void initState() {
    super.initState();

    // request permission
    Contact.requestPermission().then((granted) async {
      if (!granted) {
        Alert.message(
                "Gagal mengakses kontak. Pastikan aplikasi diizinkan untuk mengakses kontak")
            .then((value) => Screens.back());
        return;
      }

      List<Contact> contacts = await Contact.getContacts();
      setState(() {
        _granted = true;
        _contacts = contacts;
      });
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const AppBarTitle("Pilih Kontak"),
        // search button
        actions: [
          if (_contacts.isNotEmpty)
            IconButton(
              onPressed: () {
                showSearch(
                    context: context,
                    delegate: _ContactSearchDelegate(_contacts));
              },
              icon: const Icon(Icons.search),
            ),
        ],
      ),
      body: SafeArea(
        child: _buildContent(),
      ),
    );
  }

  Widget _buildSkeleton() {
    return ListTile(
      leading: Skeleton(
        width: 40,
        height: 40,
        decoration: BoxDecoration(borderRadius: BorderRadius.circular(99)),
      ),
      title: const Skeleton(
        width: 150,
        height: 15,
      ),
      subtitle: const Skeleton(
        width: 10,
        height: 10,
      ),
    );
  }

  Widget _buildContent() {
    if (!_granted) {
      return ListView.separated(
          itemBuilder: (_, __) => _buildSkeleton(),
          separatorBuilder: (_, __) => const Line(),
          itemCount: 8);
    }
    return _ContactTiles(contacts: _contacts);
  }
}

class _ContactTiles extends StatelessWidget {
  const _ContactTiles(
      {Key? key, required this.contacts, this.searchMode = false})
      : super(key: key);

  final List<Contact> contacts;
  final bool searchMode;

  @override
  Widget build(BuildContext context) {
    if (contacts.isEmpty) {
      return const Center(
        child: Text('Tidak ada kontak yang ditemukan'),
      );
    }
    return ListView.separated(
        itemBuilder: (context, index) {
          return _ContactTile(contact: contacts[index], searchMode: searchMode);
        },
        separatorBuilder: (context, index) => const Line(),
        itemCount: contacts.length);
  }
}

class _ContactTile extends StatelessWidget {
  const _ContactTile(
      {Key? key, required this.contact, required this.searchMode})
      : super(key: key);

  final Contact contact;
  final bool searchMode;

  @override
  Widget build(BuildContext context) {
    return ListTile(
      onTap: () {
        if (searchMode) {
          Navigator.pop(context); // close search
        }
        Navigator.pop(context, contact);
      },
      leading: searchMode
          ? null
          : CircleAvatar(
              backgroundColor: Theme.of(context).colorScheme.secondary,
              child: Text(contact.name[0] ?? '?',
                  style: const TextStyle(fontSize: 20)),
            ),
      title: Text(contact.name),
      subtitle: Text(contact.sanitizedPhone),
    );
  }
}

class _ContactSearchDelegate extends SearchDelegate<Contact> {
  final List<Contact> contacts;

  // make search background color same as app background color (teal)
  _ContactSearchDelegate(this.contacts)
      : super(
            searchFieldLabel: "Cari nama atau nomor telepon",
            searchFieldStyle: const TextStyle(fontSize: 17));

  @override
  List<Widget> buildActions(BuildContext context) {
    if (query.isEmpty) {
      return [];
    }
    return [
      IconButton(
        onPressed: () {
          query = '';
          showSuggestions(context);
        },
        icon: const Icon(Icons.clear),
      ),
    ];
  }

  @override
  Widget buildLeading(BuildContext context) {
    return IconButton(
      onPressed: () {
        Screens.back();
      },
      icon: const Icon(Icons.arrow_back),
    );
  }

  List<Contact> _searchContacts(String query) {
    return contacts
        .where((contact) =>
            contact.name.toLowerCase().contains(query.toLowerCase()) ||
            contact.sanitizedPhone.contains(query))
        .toList();
  }

  @override
  Widget buildResults(BuildContext context) {
    return buildSuggestions(context);
  }

  @override
  Widget buildSuggestions(BuildContext context) {
    List<Contact> contacts = _searchContacts(query);
    if (contacts.isEmpty) {
      return const Center(
        child: Text("Tidak ada hasil"),
      );
    }
    return _ContactTiles(contacts: contacts, searchMode: true);
  }
}
