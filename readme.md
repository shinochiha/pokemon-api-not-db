# Pokemon REST API
REST API ini adalah implementasi dari API Pokemon tanpa menggunakan database. REST API ini dibuat menggunakan bahasa Go dan menggunakan library `net/http` dan `github.com/gorilla/mux` untuk manajemen routing.

Pada contoh REST API Pokemon ini, informasi Pokemon disimpan sebagai variabel slice dalam struct. Slice ini berfungsi sebagai mock database, yaitu sebagai pengganti database yang sesungguhnya. Oleh karena itu, informasi Pokemon disimpan secara lokal dalam memori aplikasi dan tidak perlu terhubung ke database eksternal.

## Instalasi
Untuk menggunakan REST API ini, pertama-tama Anda harus memastikan bahwa Anda sudah menginstal Go pada komputer Anda.

1. Clone repository ini:
```bash
$ https://github.com/shinochiha/pokemon-api-not-db.git
```

2. Masuk ke direktori pokemon-api-not-db:
```bash
$ cd pokemon-api-not-db
```
3. Jalan kan perintah go mod tidy untuk mendownload package yang di perlukan
```bash
$ go mod tidy
```
4. Jalankan perintah go build untuk mengcompile file:
```bash
$ go build
```
5. Jalankan file hasil compile:
```bash
$  ./pokemon-api
```

## Penggunaan
Setelah REST API ini berjalan, Anda dapat mengaksesnya melalui browser atau aplikasi lain seperti Postman. Berikut adalah daftar endpoint yang tersedia:

#### 1. Get Pokemon
Mendapatkan daftar Pokemon. (data default nya 20)
```bash
URL: http://localhost:8000/pokemon 
Method: GET
```
Note :
jika ingin lebih dari 20 maka tambahkan params ?limit_poke_api={jumlah yang di inginkan}
```bash
Contoh: http://localhost:8000/pokemon?limit_poke_api=1279
```


#### 2. Get Detail Pokemon
Mendapatkan detail dari Pokemon tertentu.
```bash
URL: http://localhost:8000/pokemon/{id}
Method: GET
Contoh: http://localhost:8000/pokemon/1
```

#### 3. Add Pokemon
Menambahkan Pokemon baru ke daftar saya.
```bash
URL: http://localhost:8000/pokemon/{id}
Method: POST
Contoh: http://localhost:8000/pokemon/1
```
##### Note : Saat Menangkap pokemon jika success ratenya mencapai lebih dari 50% maka masukkan nickname yang di inginkan melalui terminal proses POST tidak akan selesai jika belum memasukkan nickname melalui terminal jika di bawah 50% maka akan ada pesan error `failed to catch pokemon`
#### Contoh di terminal akan muncul pesan ini :
```bash
Give a nickname for your new pokemon: mypokemonnickname
```
#### 4. My Pokemon List
Daftar Pokemon saya yang sudah di tangkap
```bash
URL: http://localhost:8000/mypokemon
Method: GET
Contoh:
{
    "count": 1,
    "results": [
        {
            "id": 1,
            "name": "bulbasaur",
            "nickname": "poke-poke",
            "height": 7,
            "weight": 69,
            "sprites": {
                "back_default": "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/back/1.png",
                "front_default": "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/1.png"
            }
        }
    ]
}
```

#### 5. Update Nickname Pokemon
Mengubah nickname data Pokemon tertentu yang ada di daftar saya.
```bash
URL: http://localhost:8000/mypokemon/{id}/{nickname}
Method: PUT
Contoh: http://localhost:8000/mypokemon/1/poke-poke
{
    "message": "Nickname changed successfully",
    "nickname": "My-bulbasaur-0"
}

nickname akan berubah menjadi ada -0 karena mengikuti fibonachi
```

#### 6. Delete Pokemon
Menghapus Pokemon tertentu dari daftar pokemon saya.
```bash
URL: http://localhost:8080/mypokemon/{id}?release={primeNumber}
Method: DELETE
Contoh: http://localhost:8000/mypokemon/1?release=5 (pilih angka prima apa saja untuk membebaskan pokemon)
{
    "message": "Pokemon has been released successfully."
}

jika memilih angka selain prima akan terkena error # Can't release the pokemon with this number
```

## Berikut penjelasan cara menggunakan REST API Sederhana ini.