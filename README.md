# yatter-backend-go
<!--
[![CircleCI](https://cci.dmm.com/gh/bootcamp-2020/yatter-backend-go.svg?style=svg)](https://cci.dmm.com/gh/bootcamp-2020/yatter-backend-go)
-->

## develop
### requirements
* docker
* GHE access rights

### start
```
git clone https://git.dmm.com/bootcamp-2020/yatter-backend-go
cd ./yatter-backend-go
docker-compose up -d
```

### modify DB model
Do NOT use any migration tools, for simplification.
`ddl/ddl.sql` will be automatically executed per `docker-compose up` (regenerate container).

Using `app/domain/object` even as DTO for DB.


## structure
[Architecture Design](https://git.dmm.com/cto-tech//wiki/Architecture-design) をベースに短期間で実装するために簡素化

```
.
├── app      ----> application core codes
│   ├── app      ----> collection of dependency injected
│   ├── config   ----> config
│   ├── domain   ----> domain layer, core business logics
│   ├── handler  ----> (interface layer & application layer), request handlers
│   └── dao      ----> (infrastructure layer), implementation of domain/repository
│
└── ddl      ----> DB definition master
```

※ 本来は usecase レイヤーが必要となるが、handler レイヤーへ合体

![モジュールの依存関係](doc/module_dependency.png)

### app
モジュールの依存関係を整理するパッケージで、DIコンテナを扱います。
今回は簡素なものになっていて、DAOの組み立てとhandlerのDAO（が提供するdomain/repository）への依存の管理のみ行っています。

### config
サーバーの設定をまとめたパッケージです。DBやlistenするポートなどの設定を取得するAPIがまとめてあります。

### domain
アプリケーションのモデルを扱うdomain層のパッケージです。

#### domain/object
ドメインに登場するデータ・モノの表現やその振る舞いを記述するパッケージです。
今回は簡単のためDTOの役割も兼ねています。

#### domain/repository
domain/objectで扱うEntityの永続化に関する処理を抽象化し、インターフェースとして定義するパッケージです。
具体的な実装はdaoパッケージで行います。

### handler
HTTPリクエストのハンドラを実装するパッケージです。
リクエストからパラメータを読み取り、エンドポイントに応じた処理を行ってレスポンスを返します。
機能の提供のために必要になる各種処理の実装は別のパッケージに切り分け、handlerは入出力に注力するケースも多いですが、今回は簡単のため統合しています。

### dao
domain/repositoryに対する実装を提供するパッケージです。
DBなど外部モジュールへアクセスし、データを永続化する処理を実装します。

## handler utilities
このテンプレートではhandler実装をサポートするユーティリティを提供しています。

### app/handler/request
リクエストの扱いに関するユーティリティをまとめています。
テンプレートにはパスパラメータ`id`の読み取りを補助する関数`IDOf`を用意しています。
```
// var r *http.Mux
r.Get("/{id}", func(w http.ResponseWriter, r *http.Request){
  id, err := request.IDOf(r)
  ...
})
```

### app/handler/httperror
エラーレスポンスを返すためのユーティリティをまとめています。
```
func SomeHandler(w http.ResponseWriter, r *http.Request) {
  ...
  if err != nil {
    httperror.InternalServerError(w, err)
	return
  }
  ...
}
```

### app/handler/auth
認証付きエンドポイントの実装のためのミドルウェア関数を提供しています。
chiの場合は`chi.Mux#Use`や`chi.Mux#With`を用いて利用できます。
- [ドキュメント](https://pkg.go.dev/github.com/go-chi/chi@v1.5.4)
- [研修資料](https://git.dmm.com/pages/bootcamp-2021/bootcamp-2021-go/yatter/http/#middleware)

ミドルウェアを埋め込んだエンドポイントでは`*http.Request`から`AccountOf`で認証と紐づくアカウントを取得できます。
```
// var r *http.Request
account := auth.AccountOf(r)
```
