export interface ConfigActionBase {
  ActionUID: number
  ActionType: string
  ActionTypeID: number
  Name: string
  BreakPoint: boolean
  TypeFeatureField: any
}

export interface Tag {
  label: string;
  len: number;
}