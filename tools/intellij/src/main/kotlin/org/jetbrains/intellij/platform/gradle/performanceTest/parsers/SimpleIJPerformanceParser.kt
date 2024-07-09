// Copyright 2000-2024 JetBrains s.r.o. and contributors. Use of this source code is governed by the Apache 2.0 license.

package org.jetbrains.intellij.platform.gradle.performanceTest.parsers

import org.gradle.api.InvalidUserDataException
import org.jetbrains.intellij.platform.gradle.models.PerformanceTestScript
import java.nio.file.Path
import java.time.Duration
import java.util.concurrent.TimeUnit
import kotlin.io.path.forEachLine

class SimpleIJPerformanceParser(private val path: Path) {

    fun parse() = PerformanceTestScript.Builder().apply {
        path.forEachLine {
            with(it) {
                when {
                    contains(Keywords.PROJECT) -> projectName(substringAfter("${Keywords.PROJECT} "))
                    contains(Keywords.ASSERT_TIMEOUT) -> assertionTimeout(substringAfter("${Keywords.ASSERT_TIMEOUT} ").convertToMillis())
                    else -> appendScriptContent(this)
                }
            }
        }
    }.build()
}

private fun String.convertToMillis() = when {
    endsWith("ms") -> removeSuffix("ms").toLong()

    endsWith("s") -> TimeUnit.MILLISECONDS.convert(
        Duration.ofSeconds(
            removeSuffix("s").toLong()
        )
    )

    endsWith("M") -> TimeUnit.MILLISECONDS.convert(
        Duration.ofMinutes(
            removeSuffix("M").toLong()
        )
    )

    endsWith("H") -> TimeUnit.MILLISECONDS.convert(
        Duration.ofHours(
            removeSuffix("H").toLong()
        )
    )

    else -> takeIf { it.isNotBlank() }?.trim()?.toLong()
} ?: throw InvalidUserDataException("Value '$this' can't be converted to milliseconds")

private class Keywords {

    companion object {
        const val PROJECT = "%%project"
        const val ASSERT_TIMEOUT = "%%assertTimeout"
    }
}
