import { APP_URL, assertVerifiableAddress, gen } from '../../../../helpers'

context('Registration', () => {
  describe('successful flow', () => {
    beforeEach(() => {
      cy.visit(APP_URL + '/auth/registration')
      cy.deleteMail()
    })

    afterEach(() => {
      cy.deleteMail()
    })

    const up = (value) => `up-${value}`
    const { email, password } = gen.identity()
    it('is able to verify the email address after sign up', () => {
      cy.register({ email, password })
      cy.login({ email, password })
      cy.session().then(assertVerifiableAddress({ isVerified: false, email }))

      cy.verifyEmail({ expect: { email } })
    })

    xit('sends the warning email on double sign up', () => {
      // FIXME https://github.com/zzpu/openuser/issues/133
      cy.clearCookies()
      cy.register({ email, password: up(password) })

      cy.verifyEmail({ expect: { email } })
    })
  })
})
