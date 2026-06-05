declare namespace Api {
  namespace Ugreen {
    interface UgreenLoginResult {
      token: string;
      refreshToken: string;
    }

    interface UgreenProfile {
      userId: string;
      userName: string;
      userType: string;
      isNew: boolean;
    }
  }
}
