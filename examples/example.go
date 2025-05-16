package example

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/googollee/sorm"
)

type User struct {
	ID           uint
	Name         string
	Email        string
	PasswordHash string
	PasswordSalt string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DisabledAt   time.Time
}

func (u *User) Schema(schema sorm.Schema) {
	schema.Field(&u.ID).Name("id").PrimaryKey()
	schema.Field(&u.Name).Name("name").NotNull().Type(sorm.VarChar(64)).UniqueIndex()
	schema.Field(&u.Email).Name("email").NotNull().Type(sorm.VarChar(128))
	schema.Field(&u.PasswordHash).Name("password_hash").NotNull().Type(sorm.VarChar(128))
	schema.Field(&u.PasswordSalt).Name("password_salt").NotNull().Type(sorm.VarChar(128))
	schema.Field(&u.CreatedAt).Name("created_at").NotNull().Type(sorm.IntegerByte(8))
	schema.Field(&u.UpdatedAt).Name("updated_at").NotNull().Type(sorm.IntegerByte(8))
	schema.Field(&u.DisabledAt).Name("disabled_at").Nullable(true).Type(sorm.IntegerByte(8))
	schema.UniqueIndex(&u.Name, &u.Email)
}

func (db *DB) CreateUser(ctx context.Context, user *User, password string) error {
	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var saltBuf [128]byte
	for i := range saltBuf {
		saltBuf[i] = byte(rand.Intn(len(letters)))
	}
	user.PasswordSalt = hex.EncodeToString(saltBuf[:])

	hash := sha256.New()
	user.PasswordHash = hex.EncodeToString(hash.Sum([]byte(password + user.PasswordSalt)))

	defer func() {
		user.PasswordSalt = ""
		user.PasswordHash = ""
	}()

	if err := sorm.Create(user).Do(ctx, db.engine).Error(); err != nil {
		return fmt.Errorf("create user %q error: %w", user.Name, err)
	}

	return nil
}

func (db *DB) ValidUser(ctx context.Context, email, password string) (id uint, err error) {
	var user struct {
		ID           uint
		PasswordHash string
		PasswordSalt string
	}

	var u User
	if e := sorm.FindFirst(&user).From(&u).Where(sorm.Eq(&u.Email, email), sorm.IsNull(&u.DisabledAt)).Do(ctx, db.engine).Error; e != nil {
		if errors.Is(e, sorm.ErrRecordNotFound) {
			err = ErrNotFound
			return
		}
		err = fmt.Errorf("find email %q error: %w", email, e)
		return
	}

	hash := sha256.New()
	want := hex.EncodeToString(hash.Sum([]byte(password + user.PasswordSalt)))

	if password != want {
		err = ErrInvalidUserOrPassword
		return
	}

	id = user.ID
	return
}

func (db *DB) DisableUser(ctx context.Context, user *User) error {
	user.DisabledAt = time.Now()
	if err := sorm.Update(user).Set(&user.DisabledAt, user.DisabledAt).Where(sorm.Eq(&user.ID, user.ID)).Do(ctx, db.engine).Error; err != nil {
		return fmt.Errorf("disable user %s error: %w", user, err)
	}

	return nil
}

type Album struct {
	ID       uint
	Path     string
	ParentID uint
}

func (a *Album) Schema(schema sorm.Schema) {
	schema.Field(&a.ID).PriaryKey()
	schema.Field(&a.Path).NotNull().Type(sorm.Text())
	schema.Field(&a.ParentID).ForeignKey(a).OnDelete()
}

func (db *DB) AlbumAndChildren(ctx context.Context, ids []uint) ([]*Album, error) {
	recursiveAlbum := "recursive_album"

	var album Album
	table := sorm.With(recursiveAlbum, sorm.Recursive()).UnionAll(
		sorm.From(&album).Where(sorm.In(&album.ID, ids)),
		sorm.From(&album).Join(
			sorm.Table(recursiveAlbum),
			sorm.Eq(&album.ParentID, sorm.Field(recursiveAlbum, "id")),
		),
	)

	// same as:
	// WITH recursive recursive_album AS (
	//   SELECT * FROM album WHERE id in [<ids>...]
	//   UNION ALL
	//   SELECT album.* FROM album JOIN recursive_album ON album.parent_id = recursive_album.ID
	// )

	var ret []*Album
	if err := sorm.FindAll(&ret).From(table).Do(ctx, db.engine); err != nil {
		return nil, fmt.Errorf("query album and children error: %w", err)
	}

	return ret, nil
}
