package cond

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_checkCond(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "пусто#1",
			args: args{
				in: "",
			},
			wantErr: true,
		},
		{
			name: "не пусто#1",
			args: args{
				in: "(eq)",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkCond(tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("checkCond() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_setIndexes(t *testing.T) {

	tests := []struct {
		name  string
		in    []rune
		want  []int
		want1 int
	}{
		{
			name:  "(eq 1 1)",
			in:    []rune("()"),
			want:  []int{1, 1},
			want1: 1,
		},
		{
			//      1   2      2 2     3      321
			name:  "(or (eq 1 1) (not  (eq 2 2)))",
			in:    []rune(`(("")(("")""))`),
			want:  []int{1, 2, -1, -1, 2, 2, 3, -2, -2, 3, -3, -3, 2, 1},
			want1: 3,
		},
		{
			//      1   2      2 2     3      321
			name:  `(or (eq 1 1) (not  (eq "(2)" "(2)")))`,
			in:    []rune(`(()(("()""()")))`),
			want:  []int{1, 2, 2, 2, 3, -1, 0, 0, -1, -2, 0, 0, -2, 3, 2, 1},
			want1: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, _ := setIndexes(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setIndexes()\n\tgot =%v!\n\twant=%v!", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("setIndexes() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestOK(t *testing.T) {

	tests := []struct {
		name    string
		in      string
		m       map[string]string
		want    bool
		wantErr bool
	}{
		{
			name: "(eq $$age$$ 22)",
			//
			//                        true
			//                  true                   true
			in:      "(eq $$age$$ 22)",
			want:    false,
			m:       map[string]string{"msisdn": "79876543210", "age": "23"},
			wantErr: false,
		},
		{
			name: "(eq $$age$$ 22)",
			//
			//                        true
			//                  true                   true
			in:      "(eq $$age$$ 22)",
			want:    true,
			m:       map[string]string{"msisdn": "79876543210", "age": "22"},
			wantErr: false,
		},
		{
			name: "(and (eq $$msisdn$$ 79876543210)  (eq $$age$$ 22))",
			//
			//                        true
			//                  true                   true
			in:      "(and (eq $$msisdn$$ 79876543210)  (eq $$age$$ 22))",
			want:    true,
			m:       map[string]string{"msisdn": "79876543210", "age": "22"},
			wantErr: false,
		},
		{
			name: "(or (eq $$msisdn$$ 79876543210)  (eq $$age$$ 22))",
			//                           true
			//                 true                   false
			in:      "(or (eq $$msisdn$$ 79876543210)  (eq $$age$$ 22))",
			m:       map[string]string{"msisdn": "79876543210", "age": "22"},
			want:    true,
			wantErr: false,
		},
		{
			name: `(or (eq $$msisdn$$ 79876543210)  (eq name "Иван \"Иванович\"  Иванов"))`,
			//                           true
			//                 true                   false
			in:      `(or (eq $$msisdn$$ 79876543210)  (eq $$name$$ "Иван \"Иванович\"  Иванов"))`,
			m:       map[string]string{"msisdn": "79876543210", "age": "22", "name": `Иван \"Иванович\"  Иванов`},
			want:    true,
			wantErr: false,
		},
		{
			name: `(or (eq $$msisdn$$ 79876543210)  (eq name "Иван \"Иванович\"  Петров"))`,
			//                           true
			//                 true                   false
			in:      `(or (eq $$msisdn$$ 79876543210)  (eq $$name$$ "Иван \"Иванович\"  Петров"))`,
			m:       map[string]string{"msisdn": "79876543210", "age": "22", "name": `Иван \"Иванович\"  Иванов`},
			want:    true,
			wantErr: false,
		},
		{
			name: `(or (eq $$msisdn$$ 79876543211)  (eq name "Иван \"Иванович\"  Иванов"))`,
			//                           true
			//                 true                   false
			in:      `(or (eq $$msisdn$$ 79876543211)  (eq $$name$$ "Иван \"Иванович\"  Иванов"))`,
			m:       map[string]string{"msisdn": "79876543210", "age": "22", "name": `Иван \"Иванович\"  Иванов`},
			want:    true,
			wantErr: false,
		},
		{
			name: `(and (eq $$msisdn$$ 79876543210)  (eq name "Иван \"Иванович\"  Петров"))`,
			//                           true
			//                 true                   false
			in:      `(and (eq $$msisdn$$ 79876543210)  (eq $$name$$ "Иван \"Иванович\"  Петров"))`,
			m:       map[string]string{"msisdn": "79876543210", "age": "22", "name": `Иван \"Иванович\"  Иванов`},
			want:    false,
			wantErr: false,
		},
		{
			name: `(and (eq $$msisdn$$ 79876543211)  (eq name "Иван \"Иванович\"  Иванов"))`,
			//                           true
			//                 true                   false
			in:      `(and (eq $$msisdn$$ 79876543211)  (eq $$name$$ "Иван \"Иванович\"  Иванов"))`,
			m:       map[string]string{"msisdn": "79876543210", "age": "22", "name": `Иван \"Иванович\"  Иванов`},
			want:    false,
			wantErr: false,
		},

		{
			name: "(and (eq $$msisdn$$ 79876543210)  (gt $$age$$ 22))",
			//
			//                        true
			//                  true                   true
			in:      "(and (eq $$msisdn$$ 79876543210)  (gt $$age$$ 22))",
			want:    false,
			m:       map[string]string{"msisdn": "79876543210", "age": "22"},
			wantErr: false,
		},
		{
			name: "(and (eq $$msisdn$$ 79876543210)  (gte $$age$$ 22))",
			//
			//                        true
			//                  true                   true
			in:      "(and (eq $$msisdn$$ 79876543210)  (gte $$age$$ 22))",
			want:    true,
			m:       map[string]string{"msisdn": "79876543210", "age": "22"},
			wantErr: false,
		},
		{
			name:    "bad syntax #1",
			in:      `(1)`,
			want:    false,
			m:       map[string]string{"1": "1"},
			wantErr: true,
		},
		{
			name: "bad syntax #2 (no placeholder)",
			in:   `(or (eq 123) (lt $$num$$ 42))`,
			want: false, // тоже самое что и true || num < 42
			// want:    false, // см. комментарий к wantErr
			m:       map[string]string{"num": "42", "f": "123"},
			wantErr: true,
			// wantErr: true, // вариант #2 - возвращать ошибку
		},
		{
			name:    "not filled placeholder",
			want:    false,
			in:      `(and (eq $$$$ 123) (lt $$num$$ 42))`,
			m:       map[string]string{"f": "123", "num": "42"},
			wantErr: true,
		},
		{
			name:    "different case operators",
			want:    false,
			in:      `(AnD (eQ $$f$$ foo) (nE $$b$$ baz))`,
			m:       map[string]string{"f": "foo", "b": "bar"},
			wantErr: true,
		},
		{
			name:    "missing fields",
			want:    false,
			in:      `(and (eq $$f$$ foo) (ne $$bar$$ bar))`,
			m:       map[string]string{"f": "foo", "b": "bar"},
			wantErr: true,
		},
		{
			name:    "num operators on strings #1",
			want:    false,
			in:      `(and (eq $$f$$ 123) (lt $$b$$ 42))`,
			m:       map[string]string{"f": "foo", "b": "bar"},
			wantErr: false,
		},
		{
			name:    "num operators on strings #2",
			want:    false,
			in:      `(and (eq $$f$$ "foo") (lt $$b$$ 42))`,
			m:       map[string]string{"f": "foo", "b": "bar"},
			wantErr: false,
		},
		{
			name:    "complex",
			want:    true,
			in:      `(or (and (eq $$f$$ foo) (eq $$b$$ bar)) (or (eq $$f$$ foobar) (eq $$b$$ bar)))`,
			m:       map[string]string{"f": "foobar", "b": "bar"},
			wantErr: false,
		}, {
			name:    "lost )",
			want:    false,
			in:      `(not (eq $$b$$ bar)`,
			m:       map[string]string{"f": "foobar", "b": "bar"},
			wantErr: true,
		},
		{
			name:    "empty cond",
			want:    false,
			in:      ``,
			m:       map[string]string{"f": "foo", "b": "bar"},
			wantErr: true,
		},
		{
			name:    "not cond",
			want:    false,
			in:      `(not (eq $$f$$ foo))`,
			m:       map[string]string{"f": "foo", "b": "bar"},
			wantErr: false,
		},
		{
			name:    "multiple not cond (and cond)",
			want:    false,
			in:      `(and (not (eq $$f$$ foo)) (not (eq $$b$$ bar)))`,
			m:       map[string]string{"f": "foo", "b": "bar"},
			wantErr: false,
		},
		{
			name:    "multiple not cond (or cond)",
			want:    true,
			in:      `(or (not (eq $$f$$ foobar)) (not (eq $$b$$ bar)))`,
			m:       map[string]string{"f": "foo", "b": "bar"},
			wantErr: false,
		},
		{
			name:    "string with braces (with quotes)",
			want:    true,
			in:      `(eq $$f$$ "string (foo)")`,
			m:       map[string]string{"f": "string (foo)"},
			wantErr: false,
		},
		{
			name:    "string with braces (with quotes)",
			want:    false,
			in:      `(eq $$f$$ "string (bar)")`,
			m:       map[string]string{"f": "string (foo)"},
			wantErr: false,
		},
		{
			name:    "string with braces (no quotes)",
			want:    false,
			in:      `(eq $$f$$ (foo))`,
			m:       map[string]string{"f": "(foo)"},
			wantErr: true,
		},
		{
			name:    "string with braces (escaped)",
			want:    false,
			in:      `(eq $$f$$ \(foo\))`,
			m:       map[string]string{"f": "(foo)"},
			wantErr: true,
		},
		{
			name:    `empty field value ("" eq "")`,
			want:    true,
			in:      `(eq $$f$$ "")`,
			m:       map[string]string{"f": ""},
			wantErr: false,
		},
		{
			name:    `empty field value`,
			want:    false,
			in:      `(eq $$f$$ foo)`,
			m:       map[string]string{"f": ""},
			wantErr: false,
		},
		{
			name:    "empty field value (complex)",
			want:    true,
			in:      `(and (eq $$f$$ "") (ne $$b$$ ""))`,
			m:       map[string]string{"f": "", "b": "123"},
			wantErr: false,
		},
		{
			name:    "eqi#1->true",
			want:    true,
			in:      `(eqi $$f$$ QWE)`,
			m:       map[string]string{"f": "qwe"},
			wantErr: false,
		},
		{
			name:    "eqi#2->true",
			want:    true,
			in:      `(eqi $$f$$ QWE)`,
			m:       map[string]string{"f": "QWE"},
			wantErr: false,
		},
		{
			name:    "eqi#1->false",
			want:    false,
			in:      `(eqi $$f$$ ASD)`,
			m:       map[string]string{"f": "qwe"},
			wantErr: false,
		},
		{
			name:    "eqi#2->false",
			want:    false,
			in:      `(eqi $$f$$ ASD)`,
			m:       map[string]string{"f": "QWE"},
			wantErr: false,
		},
		{
			name:    "contain#2->true",
			want:    true,
			in:      `(contain $$f$$ asdQWEzxc)`,
			m:       map[string]string{"f": "QWE"},
			wantErr: false,
		},
		{
			name:    "contain#2->false",
			want:    false,
			in:      `(contain $$f$$ asdqwezxc)`,
			m:       map[string]string{"f": "QWE"},
			wantErr: false,
		},
		{
			name:    "icontain#2->true",
			want:    true,
			in:      `(icontain $$f$$ asdQWEzxc)`,
			m:       map[string]string{"f": "QWE"},
			wantErr: false,
		},
		{
			name:    "icontain#2->true",
			want:    true,
			in:      `(icontain $$f$$ asdqwezxc)`,
			m:       map[string]string{"f": "QWE"},
			wantErr: false,
		},
		{
			name:    "icontain#2->false",
			want:    false,
			in:      `(icontain $$f$$ asdQWEzxc)`,
			m:       map[string]string{"f": "FGH"},
			wantErr: false,
		},
		{
			name:    "icontain#2->false",
			want:    false,
			in:      `(icontain $$f$$ asdqwezxc)`,
			m:       map[string]string{"f": "FGH"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OK(tt.in, tt.m)
			fmt.Printf("%s => got=%t, err=%v\n", tt.in, got, err)
			if (err != nil) != tt.wantErr {
				t.Errorf("OK() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("OK() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkBrackets(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want bool
	}{
		{
			name: "good#1",
			in:   "()",
			want: true,
		},
		{
			name: "good#2",
			in:   "(()())",
			want: true,
		},
		{
			name: "good#3",
			in:   "((()))",
			want: true,
		},
		{
			name: "bad#1",
			in:   "(",
			want: false,
		},
		{
			name: "bad#2",
			in:   "(()",
			want: false,
		},
		{
			name: "bad#3",
			in:   "())",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkBrackets(tt.in); got != tt.want {
				t.Errorf("checkBrackets() = %v, want %v", got, tt.want)
			}
		})
	}
}
