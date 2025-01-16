package locales

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
)

var en = &LocaleTextItems{
	DefaultTypes: &DefaultTypes{
		DecimalSeparator:    core.DECIMAL_SEPARATOR_DOT,
		DigitGroupingSymbol: core.DIGIT_GROUPING_SYMBOL_COMMA,
	},
	DataConverterTextItems: &DataConverterTextItems{
		Alipay:       "Alipay",
		WeChatWallet: "Wallet",
	},
	VerifyEmailTextItems: &VerifyEmailTextItems{
		Title:                     "Vahvista sähköposti",
		SalutationFormat:          "Hei %s,",
		DescriptionAboveBtn:       "Vahvista sähköpostiosoitteesi klikkaamalla linkkiä.",
		VerifyEmail:               "Vahvista sähköposti",
		DescriptionBelowBtnFormat: "Mikäli et avannut %s tiliä, voit jättää tämän postin huomiotta. Mikäli linkki ei toimi, kopioi ylläoleva linkki ja liitä se selaimesi osoiteriville. Vahvistuslinkki on voimassa %v minuuttia.",
	},
	ForgetPasswordMailTextItems: &ForgetPasswordMailTextItems{
		Title:                     "Nollaa salasana",
		SalutationFormat:          "Hei %s,",
		DescriptionAboveBtn:       "Saimme pyynnön nollata salasanasi. Klikkaa sllsolevaa linkkiä nollataksesi salasanan.",
		ResetPassword:             "Nollaa salasana",
		DescriptionBelowBtnFormat: "Mikäli et pyytänyt salasanan nollausta, jätä tämä posti huomiotta. Mikäli linkki ei toimi, kopioi ylläoleva linkki ja liitä se selaimesi osoiteriville. Salasanan nollauslinkki on voimassa %v minuuttia.",
	},
}
