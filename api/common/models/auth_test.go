package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_GenerateAndParseToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	id := uuid.New()

	token, err := GenerateToken(id)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	got, err := ParseToken(token)
	assert.NoError(t, err)
	assert.Equal(t, id, got, "parsed subject should equal the original user ID")
}

func Test_ParseToken_Invalid(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	_, err := ParseToken("not.a.valid.token")
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func Test_ParseToken_WrongSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "secret-a")
	token, err := GenerateToken(uuid.New())
	assert.NoError(t, err)

	// Re-parsing under a different secret must fail signature validation.
	t.Setenv("JWT_SECRET", "secret-b")
	_, err = ParseToken(token)
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func Test_Token_MissingSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "")

	_, err := GenerateToken(uuid.New())
	assert.ErrorIs(t, err, ErrMissingSecret)

	_, err = ParseToken("anything")
	assert.ErrorIs(t, err, ErrMissingSecret)
}

func Test_HashAndCheckPassword(t *testing.T) {
	const plain = "s3cr3t-pass"

	hash, err := HashPassword(plain)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, plain, hash, "hash must not equal the plaintext")

	assert.True(t, CheckPassword(hash, plain), "correct password should match")
	assert.False(t, CheckPassword(hash, "wrong-pass"), "wrong password should not match")
	assert.False(t, CheckPassword("not-a-hash", plain), "malformed hash should not match")
}
