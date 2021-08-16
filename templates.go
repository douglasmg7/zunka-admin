package main

import "html/template"

var (
	// Geral.
	tmplMaster,
	tmplIndex,
	tmplDeniedAccess,

	// Misc.
	tmplMessage,
	tmplChangelog,
	tmplTest,

	// User.
	tmplUserAdd,
	tmplUserAccount,
	tmplUserChangeName,
	tmplUserChangeEmail,
	tmplUserChangeMobile,
	tmplUserChangePassword,
	tmplUserChangeRG,
	tmplUserChangeCPF,
	tmplUserDeleteAccount,

	// Aldo.
	tmplAldoProducts,
	tmplAldoProduct,
	tmplAldoCategories,

	// Allnations.
	tmplAllnationsProducts,
	tmplAllnationsProduct,
	tmplAllnationsFilters,
	tmplAllnationsCategories,
	tmplAllnationsMakers,

	// Handytech.
	tmplHandytechProducts,
	tmplHandytechProduct,
	tmplHandytechFilters,
	tmplHandytechCategories,
	tmplHandytechMakers,

	// Mercado Livre.
	tmplMercadoLivreMessage,
	tmplMercadoLivreAuthUser,
	tmplMercadoLivreUserCode,
	tmplMercadoLivreUserInfo,
	tmplMercadoLivreActiveProducts,

	// Auth.
	tmplAuthSignup,
	tmplAuthSignin,
	tmplPasswordRecovery,
	tmplPasswordReset,

	// Student.
	tmplStudent,
	tmplAllStudent,
	tmplNewStudent *template.Template
)

// Load templates
func loadTemplates() {

	// Geral.
	tmplMaster = template.Must(template.ParseGlob("templates/master/*"))
	tmplIndex = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/index.tpl"))
	tmplDeniedAccess = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/misc/deniedAccess.tpl"))
	// Misc.
	tmplMessage = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/misc/message.tpl"))
	tmplChangelog = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/misc/changelog.gohtml"))
	tmplTest = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/misc/test.gohtml"))
	// User.
	tmplUserAdd = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/userAdd.tpl"))
	tmplUserAccount = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/user/userAccount.tpl"))
	tmplUserChangeName = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/user/userChangeName.tpl"))
	tmplUserChangeEmail = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/user/userChangeEmail.tpl"))
	tmplUserChangeMobile = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/user/userChangeMobile.tpl"))
	// Aldo.
	tmplAldoProducts = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/aldo/aldoProducts.tmpl"))
	tmplAldoProduct = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/aldo/aldoProduct.tmpl"))
	tmplAldoCategories = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/aldo/aldoCategories.tmpl"))
	// Allnations.
	tmplAllnationsProducts = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/allnations/allnationsProducts.tmpl"))
	tmplAllnationsProduct = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/allnations/allnationsProduct.tmpl"))
	tmplAllnationsFilters = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/allnations/allnationsFilters.tmpl"))
	tmplAllnationsCategories = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/allnations/allnationsCategories.tmpl"))
	tmplAllnationsMakers = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/allnations/allnationsMakers.gohtml"))
	// Handytech.
	tmplHandytechProducts = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/handytech/handytechProducts.gohtml"))
	tmplHandytechProduct = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/handytech/handytechProduct.gohtml"))
	tmplHandytechFilters = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/handytech/handytechFilters.gohtml"))
	tmplHandytechCategories = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/handytech/handytechCategories.gohtml"))
	tmplHandytechMakers = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/handytech/handytechMakers.gohtml"))
	// Mercado Livre.
	tmplMercadoLivreMessage = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/mercado_livre/mercadoLivreMessage.gohtml"))
	tmplMercadoLivreAuthUser = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/mercado_livre/mercadoLivreAuthUser.gohtml"))
	tmplMercadoLivreUserCode = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/mercado_livre/mercadoLivreUserCode.gohtml"))
	tmplMercadoLivreUserInfo = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/mercado_livre/mercadoLivreUserInfo.gohtml"))
	tmplMercadoLivreActiveProducts = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/mercado_livre/mercadoLivreActiveProducts.gohtml"))
	// tmplMercadoLivreProductDetail = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/mercado_livre/mercadoLivreProductDetail.gohtml"))

	// Auth.
	tmplAuthSignup = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/auth/signup.tpl"))
	tmplAuthSignin = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/auth/signin.tpl"))
	tmplPasswordRecovery = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/auth/passwordRecovery.tpl"))
	tmplPasswordReset = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/auth/passwordReset.tpl"))
	// Student.
	tmplStudent = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/student.tpl"))
	tmplAllStudent = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/allStudent.tpl"))
	tmplNewStudent = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/newStudent.tpl"))
}
