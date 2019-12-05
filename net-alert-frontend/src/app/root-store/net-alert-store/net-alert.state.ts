

export const selectUser = (state: NetAlertState) => state.NetworkInfo;
export const selectProfiles = (state: NetAlertState) => state.Profiles;

export interface NetAlertState{
  isLoading: boolean;
  error: any;
  NetworkInfo?: NetworkInfo;
  Owner?:Owner;
  Profiles?: Profile[];
}

export interface NetworkCard {
  Name: string; 
  Mac: string;
}
export interface Owner {
  Mac: string;
  NickName: string;
  IP: string;
  Email: string;
  Phone: string;
  GetEmailAlerts: boolean;
  GetSMSAlerts: boolean;
  LastLoginTime: Date;
}
export interface NetworkInfo {
  SSID: string;
  BSSID: string;
  GatewayIP: string;
  GatewayMAC: string;
  NetworkCards: NetworkCard[];
}



export interface Profile {
  Mac: string;
  NickName: string;
  CreateDate: Date;
  Sites?: any;
}

export interface NetworkCard {
  Name: string;
  Mac: string;
}

