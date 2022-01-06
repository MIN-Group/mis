/**
 * @Author: wzx
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/1/16 上午5:01
 *@Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package cert

//显示的去定义证书的错误，这里是x509的错误类型，选择一部分
type InvalidReason int

const (
	// Wrong version.
	IncompatibleVersion InvalidReason = iota
	// The certificate does not reach the effective date.
	NotReachEffectiveDate
	// NotAuthorizedToSign results when a certificate is signed by another
	// which isn't marked as a CA certificate.
	NotAuthorizedToSign
	// Expired results when a certificate has expired, based on the time
	// given in the VerifyOptions.
	Expired
	// Marshal or unmarshal fails.
	MarshalError
	// Error happens when trying to verify a signature
	SignatureError
	//Unknown SignatureAlgorithm
	UnknownSignatureAlgorithm


	// Below Not Use
	// CANotAuthorizedForThisName results when an intermediate or root
	// certificate has a name constraint which doesn't permit a DNS or
	// other name (including IP address) in the leaf certificate.
	CANotAuthorizedForThisName
	// TooManyIntermediates results when a path length constraint is
	// violated.
	TooManyIntermediates
	// IncompatibleUsage results when the certificate's key usage indicates
	// that it may only be used for a different purpose.

	// IncompatibleUsage results when the certificate's key usage indicates
	// that it may only be used for a different purpose.
	IncompatibleUsage
	// NameMismatch results when the subject name of a parent certificate
	// does not match the issuer name in the child.
	NameMismatch
	// NameConstraintsWithoutSANs results when a leaf certificate doesn't
	// contain a Subject Alternative Name extension, but a CA certificate
	// contains name constraints, and the Common Name can be interpreted as
	// a hostname.
	//
	// You can avoid this error by setting the experimental GODEBUG environment
	// variable to "x509ignoreCN=1", disabling Common Name matching entirely.
	// This behavior might become the default in the future.
	NameConstraintsWithoutSANs
	// UnconstrainedName results when a CA certificate contains permitted
	// name constraints, but leaf certificate contains a name of an
	// unsupported or unconstrained type.
	UnconstrainedName
	// TooManyConstraints results when the number of comparison operations
	// needed to check a certificate exceeds the limit set by
	// VerifyOptions.MaxConstraintComparisions. This limit exists to
	// prevent pathological certificates can consuming excessive amounts of
	// CPU time to verify.
	TooManyConstraints
	// CANotAuthorizedForExtKeyUsage results when an intermediate or root
	// certificate does not permit a requested extended key usage.
	CANotAuthorizedForExtKeyUsage
)

// CertificateInvalidError results when an odd error occurs. Users of this
// library probably want to handle all these errors uniformly.
type CertificateInvalidError struct {
	//Cert   *Certificate
	Reason InvalidReason
	Detail string
}

func (e CertificateInvalidError) Error() string {
	switch e.Reason {
	case IncompatibleVersion:
		return "min-cert: certificate has wrong version number"
	case NotReachEffectiveDate:
		return "min-cert: certificate does not reach the effective date: " + e.Detail
	case NotAuthorizedToSign:
		return "min-cert: certificate is not authorized to sign other certificates"
	case Expired:
		return "min-cert: certificate has expired or is not yet valid: " + e.Detail
	case MarshalError:
		return "min-cert: marshal or unmarshal fails: " + e.Detail
	case SignatureError:
		return "min-cert: Error happens when trying to verify a signature: " + e.Detail
	case UnknownSignatureAlgorithm:
		return "min-cert: Unknown SignatureAlgorithm: " + e.Detail


	// Below Not Use
	case CANotAuthorizedForThisName:
		return "min-cert: a root or intermediate certificate is not authorized to sign for this name: " + e.Detail
	case CANotAuthorizedForExtKeyUsage:
		return "min-cert: a root or intermediate certificate is not authorized for an extended key usage: " + e.Detail
	case TooManyIntermediates:
		return "min-cert: too many intermediates for path length constraint"
	case IncompatibleUsage:
		return "min-cert: certificate specifies an incompatible key usage"
	case NameMismatch:
		return "min-cert: issuer name does not match subject from issuing certificate"
	case NameConstraintsWithoutSANs:
		return "min-cert: issuer has name constraints but leaf doesn't have a SAN extension"
	case UnconstrainedName:
		return "min-cert: issuer has name constraints but leaf contains unknown or unconstrained name: " + e.Detail
	}
	return "min-cert: unknown error"
}