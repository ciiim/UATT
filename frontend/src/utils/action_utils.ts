import type { ConfigActionBase, Tag } from "../types/Action";
import { bsd_testtool } from "../../wailsjs/go/models";

export function parseActionTags(action: bsd_testtool.ConfigActionBaseJson): Tag[] {
  const tags: Tag[] = []
  const field = action.TypeFeatureField as any

  if (!field) return tags

  // ---------------------- IO 类模块 ----------------------
  if (action.ActionTypeID === 1 || action.ActionTypeID === 2) {

    // 解析子模块数组
    if (Array.isArray(field.Modules)) {
      for (const mod of field.Modules) {
        switch (mod.ModuleTypeID) {
          case 10: // Fill
            tags.push({ label: `Fill:${mod.UseVar}`, len: mod.FillLength ?? 0 })
            break
          case 11: // Fixed
            tags.push({ label: 'Fixed', len: mod.FixedContent?.length ?? 0 })
            break
          case 12: // Calc
            tags.push({ label: `Calc:${mod.CalcFunc}`, len: mod.PlaceholderBytes?.length ?? 0 })
            break
          case 13: // Custom
            tags.push({ label: 'Custom', len: mod.CustomLength ?? 0 })
            break
          default:
            tags.push({ label: 'UnknownSub', len: 0 })
        }
      }
    }
    return tags
  }

  // ---------------------- Control 类模块 ----------------------
  switch (action.ActionTypeID) {
    case 23: // Declare
      tags.push({ label: `${field.VarName}:${field.VarType}`, len: 0 })
      break
    case 24: // IF
      tags.push({ label: `IF ${field.Condition}`, len: 0 })
      break
    case 25: // ELSE
      tags.push({ label: 'ELSE', len: 0 })
      break
    case 26: // FOR
      tags.push({ label: `FOR ${field.UseVar}`, len: field.EnterCondition?.length ?? 0 })
      break
    case 27: // EndBlock
      tags.push({ label: 'EndBlock', len: 0 })
      break
    case 28: // LABEL
      tags.push({ label: `Label: ${field.LabelName}`, len: 0 })
      break
    case 29: // GOTO
      tags.push({ label: `Goto: ${field.Label}`, len: 0 })
      break
    case 30: // ChangeBaudrate
      tags.push({ label: 'BaudRate', len: field.TargetBaudRate ?? 0 })
      break
    case 31: // Stop
      tags.push({ label: 'Stop', len: field.StopCode ?? 0 })
      break
  }

  // ---------------------- Debug 类模块 ----------------------
  switch (action.ActionTypeID) {
    case 90: // Print
      tags.push({ label: `Print: ${field.PrintFmt}`, len: 0 })
      break
    case 91: // Delay
      tags.push({ label: `Delay: ${field.DelayMs}ms`, len: 0 })
      break
  }

  return tags
}

export function computeIndents(actions: any[]) {
  let indentLevel = 0
  const result = actions.map(act => {
    // 当前 action 的缩进值
    let actIndent = indentLevel

    switch (act.ActionTypeID) {
      // 块开始（IF、FOR 等）
      case 24: // IF
      case 26: // FOR
        actIndent = indentLevel
        indentLevel++ // 块开始，下一级缩进
        break
      case 25: // ELSE
        actIndent = indentLevel - 1 >= 0 ? indentLevel - 1 : 0
        break
      case 27: // EndBlock
        indentLevel = Math.max(0, indentLevel - 1)
        actIndent = indentLevel
        break
      default:
        actIndent = indentLevel
    }

    return { ...act, indent: actIndent }
  })

  return result
}